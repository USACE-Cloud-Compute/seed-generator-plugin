package main

import (
	"flag"
	"fmt"

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
}
