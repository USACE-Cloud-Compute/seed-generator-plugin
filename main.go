package seedgenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/usace/seed-generator/seedgeneratormodel"
	"github.com/usace/wat-go"
)

func main() {
	fmt.Println("event generator!")
	pm, err := wat.InitPluginManager()
	if err != nil {
		pm.LogMessage(wat.Message{
			Message: err.Error(),
		})
	}
	payload, err := pm.GetPayload()
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return
	}
	err = computePayload(payload, pm)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return
	}
}
func computePayload(payload wat.Payload, pm wat.PluginManager) error {
	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	var modelResourceInfo wat.DataSource
	found := false
	for _, rfd := range payload.Inputs {
		if strings.Contains(rfd.Name, "seedgenerator.json") {
			modelResourceInfo = rfd
			found = true
		}
	}
	if !found {
		err := fmt.Errorf("could not find %s.json", "seedgenerator")
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	modelBytes, err := pm.GetObject(modelResourceInfo)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	var eventGeneratorModel seedgeneratormodel.Model
	err = json.Unmarshal(modelBytes, &eventGeneratorModel)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	eventIndex := pm.EventNumber()
	modelResult, err := eventGeneratorModel.Compute(eventIndex)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	err = pm.PutObject(payload.Outputs[0], bytes)
	if err != nil {
		pm.LogError(wat.Error{
			ErrorLevel: wat.ERROR,
			Error:      err.Error(),
		})
		return err
	}
	pm.ReportProgress(wat.StatusReport{
		Status:   wat.SUCCEEDED,
		Progress: 100,
	})
	return nil
}
