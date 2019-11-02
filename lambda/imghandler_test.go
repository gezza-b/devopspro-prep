package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
)

const expectedResp = "Hello"
const putObjectEvent = "../testevents/s3upload.json"

func TestUploadEvent(t *testing.T) {
	file, _ := ioutil.ReadFile(putObjectEvent)

	data := MyEvent{}
	_ = json.Unmarshal([]byte(file), &data)
	var err, res = HandleRequest(data)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, expectedResp)
}
