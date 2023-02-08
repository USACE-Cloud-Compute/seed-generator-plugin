package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/usace/cc-go-sdk"
	"github.com/usace/seed-generator/seedgeneratormodel"
)

func main() {
	fmt.Println("event generator!")
	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}
	payload := pm.GetPayload()
	err = computePayload(payload, pm)
	if err != nil {
		pm.LogError(cc.Error{
			ErrorLevel: cc.ERROR,
			Error:      err.Error(),
		})
	}
}

func computePayload(payload cc.Payload, pm *cc.PluginManager) error {
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

	var eventGeneratorModel seedgeneratormodel.Model
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

	outds, err := pm.GetOutputDataSource("seedoutput")
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
	pm.ReportProgress(cc.StatusReport{
		Status:   cc.SUCCEEDED,
		Progress: 100,
	})
	return nil
}
