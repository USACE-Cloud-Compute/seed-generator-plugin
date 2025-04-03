package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/usace/cc-go-sdk"
	tiledb "github.com/usace/cc-go-sdk/tiledb-store"
	"github.com/usace/seed-generator/blockgeneratormodel"
	"github.com/usace/seed-generator/seedgeneratormodel"
)

var logContext = context.Background()

func main() {
	fmt.Println("seed generator!")
	//register tiledb
	cc.DataStoreTypeRegistry.Register("TILEDB", tiledb.TileDbEventStore{})
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	payload := pm.Payload
	err = computePayloadActions(payload, pm)
	if err != nil {
		pm.Logger.Log(logContext, slog.LevelError, err.Error())
		return
	}
	pm.Logger.SendMessage("whatchannel?", "compute complete", slog.Attr{Key: "progress", Value: slog.IntValue(100)})
}

func computePayloadActions(payload cc.Payload, pm *cc.PluginManager) error {
	for _, action := range payload.Actions {
		switch action.Type {
		case "block_generation":
			err := generateBlocks(action, pm)
			if err != nil {
				return err
			}
		case "realization_seed_generation":
			err := generateSeeds(payload)
			if err != nil {
				return err
			}
		case "block_event_seed_generation":
			err := generateEventSeedsFromBlocks(action, pm)
			if err != nil {
				return err
			}
		case "block_all_seed_generation":
			err := generateAllSeedsFromBlocks(action, pm)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("could not process action of type %s", action.Type)
		}
	}

	return nil
}

func generateBlocks(action cc.Action, pm *cc.PluginManager) error {
	//initialize a blockgeneratormodel
	blockGenerator := blockgeneratormodel.BlockGenerator{
		TargetTotalEvents:    action.Attributes.GetInt64OrFail("target_total_events"),
		BlocksPerRealization: action.Attributes.GetIntOrFail("blocks_per_realization"),
		TargetEventsPerBlock: action.Attributes.GetIntOrFail("target_events_per_block"),
		Seed:                 action.Attributes.GetInt64OrDefault("seed", 1234),
	}
	blocks := blockGenerator.GenerateBlocks()
	bytedata, err := json.Marshal(blocks)
	if err != nil {
		return fmt.Errorf("could not encode blocks, %v", err)
	}
	outputDataset_name := action.Attributes.GetStringOrFail("outputdataset_name")
	storeType := action.Attributes.GetStringOrFail("store_type")
	if storeType == "eventstore" {
		recordset, err := cc.NewEventStoreRecordset(pm, &blocks, "eventstore", outputDataset_name)
		if err != nil {
			log.Fatal(err)
		}
		err = recordset.Create()
		if err != nil {
			log.Fatal(err)
		}
		err = recordset.Write(&blocks)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		breader := bytes.NewReader(bytedata)
		_, err = pm.IOManager.Put(cc.PutOpInput{
			SrcReader: breader,
			DataSourceOpInput: cc.DataSourceOpInput{
				DataSourceName: outputDataset_name,
				PathKey:        "default",
			},
		})
	}

	if err != nil {
		return fmt.Errorf("could not write file error: %v", err)
	}
	return nil
}
func generateEventSeedsFromBlocks(action cc.Action, pm *cc.PluginManager) error {

	var eventGeneratorModel seedgeneratormodel.BlockModel
	eventGeneratorModel.InitialEventSeed = action.Attributes.GetInt64OrFail("initial_event_seed")
	eventGeneratorModel.InitialRealizationSeed = action.Attributes.GetInt64OrFail("initial_realization_seed")
	plugins, err := action.Attributes.GetStringSlice("plugins")
	if err != nil {
		return err
	}
	eventGeneratorModel.Plugins = plugins
	blockreader, err := action.IOManager.GetReader(cc.DataSourceOpInput{
		DataSourceName: "blockfile",
		PathKey:        "default",
	})
	if err != nil {
		return err
	}
	defer blockreader.Close()
	var blocks []blockgeneratormodel.Block
	err = json.NewDecoder(blockreader).Decode(&blocks)
	if err != nil {
		return err
	}

	modelResult, err := eventGeneratorModel.Compute(pm.EventNumber(), blocks)
	if err != nil {
		return err
	}

	bytedata, err := json.Marshal(modelResult)
	if err != nil {
		return err
	}
	breader := bytes.NewReader(bytedata)
	_, err = action.Put(cc.PutOpInput{
		SrcReader: breader,
		DataSourceOpInput: cc.DataSourceOpInput{
			DataSourceName: "seeds",
			PathKey:        "default",
		},
	})
	if err != nil {
		return err
	}
	return nil
}
func generateAllSeedsFromBlocks(action cc.Action, pm *cc.PluginManager) error {

	var eventGeneratorModel seedgeneratormodel.BlockModel
	eventGeneratorModel.InitialEventSeed = action.Attributes.GetInt64OrFail("initial_event_seed")
	eventGeneratorModel.InitialRealizationSeed = action.Attributes.GetInt64OrFail("initial_realization_seed")
	plugins, err := action.Attributes.GetStringSlice("plugins")

	if err != nil {
		return err
	}
	eventGeneratorModel.Plugins = plugins
	var blocks []blockgeneratormodel.Block
	block_name := action.Attributes.GetStringOrFail("block_dataset_name")
	blockStoreType := action.Attributes.GetStringOrFail("block_store_type")
	seed_name := action.Attributes.GetStringOrFail("seed_dataset_name")
	seedStoreType := action.Attributes.GetStringOrFail("seed_store_type")
	if blockStoreType == "eventstore" {
		recordset, err := cc.NewEventStoreRecordset(pm, &blocks, "eventstore", block_name)
		if err != nil {
			log.Fatal(err)
		}
		//err = recordset.Create()
		if err != nil {
			log.Fatal(err)
		}
		arrayResult, err := recordset.Read()
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < arrayResult.Size(); i++ {
			var block blockgeneratormodel.Block
			arrayResult.Scan(&block)
			blocks = append(blocks, block)
		}
	} else {

		blockreader, err := action.IOManager.GetReader(cc.DataSourceOpInput{
			DataSourceName: block_name,
			PathKey:        "default",
		})
		if err != nil {
			return err
		}
		defer blockreader.Close()

		err = json.NewDecoder(blockreader).Decode(&blocks)
		if err != nil {
			return err
		}
	}

	modelResult, err := eventGeneratorModel.ComputeAll(blocks)
	if err != nil {
		return err
	}

	if seedStoreType == "eventstore" {
		//get the store
		ds, err := pm.GetStore(seedStoreType)
		if err != nil {
			return err
		}
		//check if it is the right store type
		tdbds, ok := ds.Session.(cc.MultiDimensionalArrayStore)
		if ok {
			arrayInput := cc.CreateArrayInput{
				ArrayPath: seed_name,
				Dimensions: []cc.ArrayDimension{
					{
						Name:          "Plugins",
						DimensionType: cc.DIMENSION_INT,
						Domain:        []int64{1, int64(len(plugins))},
						TileExtent:    int64(len(plugins)),
					}, {
						Name:          "Events",
						DimensionType: cc.DIMENSION_INT,
						Domain:        []int64{1, int64(len(modelResult))},
						TileExtent:    1,
					}, /*{
						Name: "Seeds",
						DimensionType: cc.DIMENSION_INT,
						Domain: []int64{1,2},
						TileExtent: 1,//what does this mean really.
					}*/
				},
				Attributes: []cc.ArrayAttribute{
					{Name: "realization_seed", DataType: cc.ATTR_INT64},
					{Name: "event_seed", DataType: cc.ATTR_INT64},
				},
				ArrayType:  cc.ARRAY_DENSE,
				CellLayout: cc.ROWMAJOR,
				TileLayout: cc.COLMAJOR,
			}
			err = tdbds.CreateArray(arrayInput)
			if err != nil {
				return err
			}

			//now make a put array input and put the data properly arranged.
			eventseeddata := make([]int64, (len(plugins))*(len(modelResult)))
			realizationseeddata := make([]int64, len(plugins)*len(modelResult))
			for i, ec := range modelResult {
				for j, plugin := range plugins {
					ss := ec.Seeds[plugin]
					eventseeddata[(i*len(plugins))+j] = ss.EventSeed
					realizationseeddata[(i*len(plugins))+j] = ss.RealizationSeed
				}
			}
			//create a buffer
			buffer := []cc.PutArrayBuffer{
				{
					AttrName: "event_seed",
					Buffer:   eventseeddata,
				},
				{
					AttrName: "realization_seed",
					Buffer:   realizationseeddata,
				},
			}
			//create an input
			input := cc.PutArrayInput{
				Buffers:   buffer,
				DataPath:  seed_name,
				ArrayType: cc.ARRAY_DENSE,
			}
			err = tdbds.PutArray(input)
			if err != nil {
				return err
			}
			mds, ok := ds.Session.(cc.MetadataStore)
			if ok {
				err = mds.PutMetadata("seed_columns", plugins)
				if err != nil {
					return err
				}

			} else {
				return errors.New("could not store metadata on which plugins recieve seeds")
			}
		} else {
			//store does not support this data.
			return errors.New("store session is not a multidimensional array store")
		}
	} else {
		bytedata, err := json.Marshal(modelResult)
		if err != nil {
			return err
		}
		breader := bytes.NewReader(bytedata)
		_, err = action.Put(cc.PutOpInput{
			SrcReader: breader,
			DataSourceOpInput: cc.DataSourceOpInput{
				DataSourceName: seed_name,
				PathKey:        "default",
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func generateSeeds(payload cc.Payload) error {
	pm, err := cc.InitPluginManager()
	if err != nil {
		return fmt.Errorf("Unable to initialize the plugin manager: %s\n", err)
	}
	if len(payload.Outputs) != 1 {
		return errors.New("more than one output was defined")
	}
	reader, err := pm.IOManager.GetReader(cc.DataSourceOpInput{
		DataSourceName: "seedgenerator",
		PathKey:        "default",
	})
	if err != nil {
		return err
	}
	defer reader.Close()

	var eventGeneratorModel seedgeneratormodel.RealizationModel
	err = json.NewDecoder(reader).Decode(&eventGeneratorModel)
	if err != nil {
		return err
	}

	eventIndex := pm.EventNumber()
	modelResult, err := eventGeneratorModel.Compute(eventIndex)
	if err != nil {
		return err
	}

	bytedata, err := json.Marshal(modelResult)
	if err != nil {
		return err
	}
	breader := bytes.NewReader(bytedata)

	_, err = pm.IOManager.Put(cc.PutOpInput{
		SrcReader: breader,
		DataSourceOpInput: cc.DataSourceOpInput{
			DataSourceName: "seeds",
			PathKey:        "default",
		},
	})
	if err != nil {
		return err
	}
	return nil
}
