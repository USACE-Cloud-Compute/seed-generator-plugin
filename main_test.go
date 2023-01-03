package main

/*
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
	err = plugin.InitConfigFromEnv()
	if err != nil {
		logError(err, plugin.ModelPayload{Id: "unknownpayloadid"})
		return
	}
	mp, err = plugin.LoadPayload(path)
	if err != nil {
		logError(err, plugin.ModelPayload{Id: "unknownpayloadid"})
		return
	}

	err = computePayload(mp)
	if err != nil {
		t.Fail()
	}
}
*/
