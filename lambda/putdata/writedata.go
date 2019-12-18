package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(name MyMetaData) (response string, err error) {
	fmt.Println(":writedata:handler")
	return
	// return nil, nil
}

type MyMetaData struct {
	Name string `json:"name"`
}

func main() {
	lambda.Start(HandleRequest)
}
