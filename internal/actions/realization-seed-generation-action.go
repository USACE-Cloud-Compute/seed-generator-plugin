package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/usace-cloud-compute/cc-go-sdk"
	. "github.com/usace-cloud-compute/seed-generator/internal/compute"
)

const (
	realizationSeedGenerationActionName string = "realization-seed-generation"
)

func init() {
	cc.ActionRegistry.RegisterAction(realizationSeedGenerationActionName, &RealizationSeedGenerationAction{})
}

type RealizationSeedGenerationAction struct {
	cc.ActionRunnerBase
}

func (rsga *RealizationSeedGenerationAction) Run() error {
	pm := rsga.PluginManager
	payload := pm.Payload

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

	var eventGeneratorModel RealizationModel
	err = json.NewDecoder(reader).Decode(&eventGeneratorModel)
	if err != nil {
		return err
	}

	eventIndex, err := strconv.Atoi(pm.EventIdentifier)
	if err != nil {
		return err
	}
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
