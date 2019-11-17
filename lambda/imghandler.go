package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

const resp string = "Hello"
const bucket string = "imghandler-gez"
const region string = "us-east-1"
const imgName string = "2faces.jpeg"

type MyEvent struct {
	Name string `json:"name"`
}

// https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-context.html
// func HandleRequest(ctx context.Context, name MyEvent) (response string, err error) {
func HandleRequest(name MyEvent) (response string, err error) {
	// fmt.Println("::ctx: ", ctx)
	// lc, _ := lambdacontext.FromContext(ctx)
	// reqId := lc.AwsRequestID
	// fmt.Println(":reqID: ", reqId)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Println("NewSession error: ", err)
		return "SendMsg error while creating session", err
	}

	client := rekognition.New(sess)

	input := &rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(imgName),
			},
		},
		MaxLabels:     aws.Int64(100),
		MinConfidence: aws.Float64(70.000000),
	}
	result, err := client.DetectLabels(input)

	if err != nil {
		fmt.Println("error")
	} else {
		//https://blog.golang.org/json-and-go
		var f interface{}
		err := json.Unmarshal(b, &f)

		fmt.Println(":result: ", result)
	}
	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
