package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/usace/cc-go-sdk"
	"github.com/usace/seed-generator/blockgeneratormodel"
	"github.com/usace/seed-generator/seedgeneratormodel"
)

func main() {
	fmt.Println("seed generator!")
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	payload := pm.GetPayload()
	err = computePayloadActions(payload)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return
	}
	pm.ReportProgress(cc.StatusReport{
		Status:   cc.SUCCEEDED,
		Progress: 100,
	})
}

func computePayloadActions(payload cc.Payload) error {
	for _, action := range payload.Actions {
		switch action.Name {
		case "block_generation":
			generateBlocks(action)
		case "realization_seed_generation":
			generateSeeds(payload)
		case "block_seed_generation":
			generateSeedsFromBlocks(payload)
		default:
			log.Fatalf("%s.\n", action.Name)
		}
	}

	return nil
}

func generateBlocks(action cc.Action) {
	//initialize a blockgeneratormodel
	blockGenerator := blockgeneratormodel.BlockGenerator{
		TargetTotalEvents:    action.Parameters.GetInt64OrFail("target_total_events"),
		BlocksPerRealization: action.Parameters.GetIntOrFail("blocks_per_realization"),
		TargetEventsPerBlock: action.Parameters.GetIntOrFail("target_events_per_block"),
		Seed:                 action.Parameters.GetInt64OrDefault("seed", 1234),
	}
	blocks := blockGenerator.GenerateBlocks()
	bytes, err := json.Marshal(blocks)
	if err != nil {
		log.Fatal("could not encode blocks")
	}
	outputDataset_name := action.Parameters.GetStringOrFail("outputdataset_name")
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatal("could not init plugin manager")
	}
	outputDataset, err := pm.GetOutputDataSource(outputDataset_name)
	if err != nil {
		log.Fatal("could not find datasource")
	}
	err = pm.PutFile(bytes, outputDataset, 0)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not write file error: %v", err))
	}
}
func generateSeedsFromBlocks(payload cc.Payload) error {
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	reader, err := pm.FileReaderByName("seedgeneratorconfig", 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	defer reader.Close()

	var eventGeneratorModel seedgeneratormodel.BlockModel
	err = json.NewDecoder(reader).Decode(&eventGeneratorModel)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	blockreader, err := pm.FileReaderByName("blockfile", 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	defer blockreader.Close()
	var blocks []blockgeneratormodel.Block
	err = json.NewDecoder(reader).Decode(&blocks)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	eventIndex := pm.EventNumber()
	modelResult, err := eventGeneratorModel.Compute(eventIndex, blocks)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	bytes, err := json.Marshal(modelResult)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	outds, err := pm.GetOutputDataSource("seeds")
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	err = pm.PutFile(bytes, outds, 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	return nil
}

func generateSeeds(payload cc.Payload) error {
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	reader, err := pm.FileReaderByName("seedgenerator", 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	defer reader.Close()

	var eventGeneratorModel seedgeneratormodel.RealizationModel
	err = json.NewDecoder(reader).Decode(&eventGeneratorModel)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	eventIndex := pm.EventNumber()
	modelResult, err := eventGeneratorModel.Compute(eventIndex)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	bytes, err := json.Marshal(modelResult)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	outds, err := pm.GetOutputDataSource("seeds")
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}

	err = pm.PutFile(bytes, outds, 0)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	return nil
}
