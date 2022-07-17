package main

import (
	"encoding/json"
	"flag"
	"fmt"

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
	err := plugin.InitConfigFromPath("howDoIGetThePluginConfigPassedIn?")
	if err != nil {
		logError(err, plugin.ModelPayload{Id: "unknownpayloadid"})
		return
	}
	payload, err := plugin.LoadPayload(payloadPath)
	if err != nil {
		logError(err, plugin.ModelPayload{Id: "unknownpayloadid"})
		return
	}
	//how am i supposed to know which file is the input model file?
	//i dont know what the model name is - i need to know that to find the right file to load.
	modelResourceInfo := payload.Inputs[0].ResourceInfo //wag
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
