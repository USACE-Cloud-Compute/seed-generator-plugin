package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/henrygeorgist/eventgenerator/eventgenerator"
	"github.com/usace/wat-go-sdk/plugin"
)

func main() {
	fmt.Println("event generator!")
	var payloadPath string
	flag.StringVar(&payloadPath, "payload", "", "please specify an input file using `-payload=pathtopayload.yml`")
	flag.Parse()
	if payloadPath == "" {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   "given a blank path...\n\tplease specify an input file using `-payload=pathtopayload.yml`",
			Sender:    "eventgenerator",
			PayloadId: "unknown payloadid because the plugin package could not be properly initalized",
		})
		return
	}
	err := plugin.InitConfigFromEnv()
	if err != nil {
		logError(err, plugin.ModelPayload{Id: "unknownpayloadid"})
		return
	}
	payload, err := plugin.LoadPayload(payloadPath)
	if err != nil {
		logError(err, plugin.ModelPayload{Id: "unknownpayloadid"})
		return
	}
	err = computePayload(payload)
	if err != nil {
		logError(err, payload)
		return
	}
}
func computePayload(payload plugin.ModelPayload) error {
	if len(payload.Outputs) != 1 {
		err := errors.New("more than one output was defined")
		logError(err, payload)
		return err
	}
	var modelResourceInfo plugin.ResourceInfo
	found := false
	for _, rfd := range payload.Inputs {
		if strings.Contains(rfd.FileName, payload.Model.Name+".json") {
			modelResourceInfo = rfd.ResourceInfo
			found = true
		}
	}
	if !found {
		err := fmt.Errorf("could not find %s.json", payload.Model.Name)
		logError(err, payload)
		return err
	}
	modelBytes, err := plugin.DownloadObject(modelResourceInfo)
	if err != nil {
		logError(err, payload)
		return err
	}
	var eventGeneratorModel eventgenerator.Model
	err = json.Unmarshal(modelBytes, &eventGeneratorModel)
	if err != nil {
		logError(err, payload)
		return err
	}
	modelResult, err := eventGeneratorModel.Compute(payload.EventIndex)
	if err != nil {
		logError(err, payload)
		return err
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		logError(err, payload)
		return err
	}
	err = plugin.UpLoadFile(payload.Outputs[0].ResourceInfo, bytes)
	if err != nil {
		logError(err, payload)
		return err
	}
	plugin.Log(plugin.Message{
		Status:    plugin.SUCCEEDED,
		Progress:  100,
		Level:     plugin.INFO,
		Message:   "event generation complete",
		Sender:    "eventgenerator",
		PayloadId: payload.Id,
	})
	return nil
}
func logError(err error, payload plugin.ModelPayload) {
	plugin.Log(plugin.Message{
		Status:    plugin.FAILED,
		Progress:  0,
		Level:     plugin.ERROR,
		Message:   err.Error(),
		Sender:    "eventgenerator",
		PayloadId: payload.Id,
	})
}
