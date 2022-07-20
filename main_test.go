package main

import (
	"os"
	"testing"

	"github.com/usace/wat-go-sdk/plugin"
	"gopkg.in/yaml.v3"
)

func TestReadPayload(t *testing.T) {
	path := "./exampledata/payload.yaml"
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fail()
	}
	mp := plugin.ModelPayload{}
	err = yaml.Unmarshal(b, &mp)
	if err != nil {
		t.Fail()
	}
}
func TestComputePayload(t *testing.T) {
	path := "./exampledata/payload.yaml"
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fail()
	}
	mp := plugin.ModelPayload{}
	err = yaml.Unmarshal(b, &mp)
	if err != nil {
		t.Fail()
	}
	err = computePayload(mp)
	if err != nil {
		t.Fail()
	}
}
