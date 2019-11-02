package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	// "github.com/aws/aws-xray-sdk-go/xray"
)

const resp string = "Hello"

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(name MyEvent) (response string, err error) {
	// xray.Configure(xray.Config{
	// 	LogLevel:       "info", // default
	// 	ServiceVersion: "1.2.3",
	// })

	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
