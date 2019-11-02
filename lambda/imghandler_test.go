package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

const expectedRes = "Hello"
const putObjectEvent = "../testevents/s3upload.json"

func TestSpanish(t *testing.T) {
	file, _ := ioutil.ReadFile(putObjectEvent)
	fmt.Println("file: ", file)
}
