package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedLength = len("9476cd6b-7f9e-53b6-9942-9b1232f02810")

const putObjectEvent = "../../testevents/s3upload.json"

func TestUploadEvent(t *testing.T) {
	file, _ := ioutil.ReadFile(putObjectEvent)
	assert.NotEmpty(t, file)
	data := MyEvent{}

	_ = json.Unmarshal([]byte(file), &data)

	// var res, err = HandleRequest(context.TODO(), data)
	var res, err = HandleRequest(nil, data)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(res), expectedLength)
}
