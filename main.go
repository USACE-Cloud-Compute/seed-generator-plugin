package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println("event generator!")
	var payload string
	flag.StringVar(&payload, "payload", "", "please specify an input file using `-payload=pathtopayload.yml`")
	flag.Parse()

	if payload == "" {
		fmt.Println("given a blank path...")
		fmt.Println("please specify an input file using `-payload=pathtopayload.yml`")
		return
	}
	payloadInstructions, err := utils.LoadModelPayloadFromS3(payload, fs)
	if err != nil {
		fmt.Println("not successful", err)
		return
	}

	message := fmt.Sprintf("Fragility Curve Complete %v", val)
	fmt.Println("sending message: " + message)
	queueURL := fmt.Sprintf("%v/queue/events", queue.Endpoint)
	fmt.Println("sending message to:", queueURL)
	_, err = queue.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     &queueURL,
	})
	key := payloadInstructions.PluginImageAndTag + "_" + payloadInstructions.Name + "_R" + fmt.Sprint(payloadInstructions.Realization.Index) + "_E" + fmt.Sprint(payloadInstructions.Event.Index)
	out := red.Set(key, "complete", 0)
	fmt.Println(out)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(message)
	return
}
