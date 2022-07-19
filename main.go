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
	if len(payload.Outputs) != 1 {
		logError(errors.New("more than one output was defined"), payload)
		return
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
		logError(fmt.Errorf("could not find %s.json", payload.Model.Name), payload)
		return
	}
	modelBytes, err := plugin.DownloadObject(modelResourceInfo)
	if err != nil {
		logError(err, payload)
		return
	}
	var eventGeneratorModel eventgenerator.Model
	err = json.Unmarshal(modelBytes, &eventGeneratorModel)
	if err != nil {
		logError(err, payload)
		return
	}
	modelResult, err := eventGeneratorModel.Compute(payload.EventIndex)
	if err != nil {
		logError(err, payload)
		return
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		logError(err, payload)
		return
	}
	err = plugin.UpLoadFile(payload.Outputs[0].ResourceInfo, bytes)
	if err != nil {
		logError(err, payload)
		return
	}
	plugin.Log(plugin.Message{
		Status:    plugin.SUCCEEDED,
		Progress:  100,
		Level:     plugin.INFO,
		Message:   "event generation complete",
		Sender:    "eventgenerator",
		PayloadId: payload.Id,
	})
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
