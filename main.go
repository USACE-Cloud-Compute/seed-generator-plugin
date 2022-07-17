package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/henrygeorgist/eventgenerator/model"
	"github.com/usace/wat-go-sdk/plugin"
)

func main() {
	fmt.Println("event generator!")
	var payloadPath string
	var localWorkingDir string = "/workingdir"
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
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: "unknown payloadid because the plugin package could not be properly initalized",
		})
		return
	}
	payload, err := plugin.LoadPayload(payloadPath)
	if err != nil {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: "unknown payloadid because the payload did not load",
		})
		return
	}
	plugin.Log(plugin.Message{
		Status:    plugin.SUCCEEDED,
		Progress:  100,
		Level:     plugin.INFO,
		Message:   "awwww yeah",
		Sender:    "eventgenerator",
		PayloadId: payload.Id,
	})

	err = plugin.CopyPayloadInputsLocally(payload, localWorkingDir)
	if err != nil {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: payload.Id,
		})
		return
	}
	//how am i supposed to know which file is the input model file?
	//i dont know what the model name is - i need to know that to find the right file to load.
	modelResourceInfo := payload.Inputs[0].ResourceInfo //wag
	eventGeneratorModel, err := model.Init(modelResourceInfo)
	if err != nil {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: payload.Id,
		})
		return
	}
	modelResult, err := eventGeneratorModel.Compute(payload.EventIndex)
	if err != nil {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: payload.Id,
		})
		return
	}
	bytes, err := json.Marshal(modelResult)
	if err != nil {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: payload.Id,
		})
		return
	}
	err = plugin.UpLoadFile(payload.Outputs[0].ResourceInfo, bytes)
	if err != nil {
		plugin.Log(plugin.Message{
			Status:    plugin.FAILED,
			Progress:  0,
			Level:     plugin.ERROR,
			Message:   err.Error(),
			Sender:    "eventgenerator",
			PayloadId: payload.Id,
		})
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
