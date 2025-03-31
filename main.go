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
	err = computePayloadActions(payload)
	if err != nil {
		pm.Logger.Log(logContext, slog.LevelError, err.Error())
		return
	}
	pm.Logger.SendMessage("whatchannel?", "compute complete", slog.Attr{Key: "progress", Value: slog.IntValue(100)})
}

func computePayloadActions(payload cc.Payload) error {
	for _, action := range payload.Actions {
		switch action.Type {
		case "block_generation":
			err := generateBlocks(action)
			if err != nil {
				return err
			}
		case "realization_seed_generation":
			err := generateSeeds(payload)
			if err != nil {
				return err
			}
		case "block_seed_generation":
			err := generateSeedsFromBlocks(payload)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("could not process action of type %s", action.Type)
		}
	}

	return nil
}

func generateBlocks(action cc.Action) error {
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
	pm, err := cc.InitPluginManager()
	if err != nil {
		return fmt.Errorf("could not init plugin manager, %v", err)
	}
	breader := bytes.NewReader(bytedata)
	_, err = pm.IOManager.Put(cc.PutOpInput{
		SrcReader: breader,
		DataSourceOpInput: cc.DataSourceOpInput{
			DataSourceName: outputDataset_name,
			PathKey:        "default",
		},
	})
	if err != nil {
		return fmt.Errorf("could not write file error: %v", err)
	}
	return nil
}
func generateSeedsFromBlocks(payload cc.Payload) error {
	pm, err := cc.InitPluginManager()
	if err != nil {
		return fmt.Errorf("Unable to initialize the plugin manager: %s\n", err)
	}
	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		return err
	}
	reader, err := pm.IOManager.GetReader(cc.DataSourceOpInput{
		DataSourceName: "seedgeneratorconfig",
		PathKey:        "default",
	})
	if err != nil {
		return err
	}
	defer reader.Close()

	var eventGeneratorModel seedgeneratormodel.BlockModel
	err = json.NewDecoder(reader).Decode(&eventGeneratorModel)
	if err != nil {
		return err
	}
	blockreader, err := pm.IOManager.GetReader(cc.DataSourceOpInput{
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
	eventIndex := pm.EventNumber()
	modelResult, err := eventGeneratorModel.Compute(eventIndex, blocks)
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
