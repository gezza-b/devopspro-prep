package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedResp = "Hello"
const putObjectEvent = "../testevents/s3upload.json"

func TestUploadEvent(t *testing.T) {
	file, _ := ioutil.ReadFile(putObjectEvent)
	assert.NotEmpty(t, file)
	data := MyEvent{}

	_ = json.Unmarshal([]byte(file), &data)
	// var res, err = HandleRequest(context.TODO(), data)
	var res, err = HandleRequest(data)
	assert.Equal(t, err, nil)
	assert.Equal(t, res, expectedResp)
}
