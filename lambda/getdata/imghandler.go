package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/sns"
)

const resp string = "Hello"
const bucket string = "imghandler-gez"
const imgName string = "2faces.jpeg"
const maxlabels = 100
const minConfidence float64 = 75.000000

var region string = os.Getenv("Region")
var account string = os.Getenv("Account")

func init() {
	if region == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		region = os.Getenv("Region")
		account = os.Getenv("Account")
	}
}

type MyEvent struct {
	Name string `json:"name"`
}

// https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-context.html
func HandleRequest(name MyEvent) (response string, err error) {
	// lc, _ := lambdacontext.FromContext(ctx)
	// reqId := lc.AwsRequestID
	fmt.Println(":event: ", name)

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
		MaxLabels:     aws.Int64(maxlabels),
		MinConfidence: aws.Float64(minConfidence),
	}

	resultJson, err := client.DetectLabels(input)
	if err != nil {
		fmt.Println("error: ", err)
	}

	var info = getImgInfo(resultJson)
	fmt.Println(":info: ", info)

	// sendSns(jsonData)
	return resp, nil
}

func deduplicate(duplicates []string) (dedupe []string) {
	check := make(map[string]int)
	for _, val := range duplicates {
		check[val] = 1
	}

	for item, _ := range check {
		dedupe = append(dedupe, item)
	}
	return dedupe
}

func getImgInfo(res *rekognition.DetectLabelsOutput) (imgInfo ImgInfo) {
	labels := res.Labels
	parents := make([]string, 0)
	objects := make([]string, 0)

	for i := range labels {
		// needs to meet minimum confidence
		if *labels[i].Confidence > minConfidence {
			var l rekognition.Label = *labels[i]
			// parents
			for i := 0; i < len(l.Parents); i++ {
				parents = append(parents, *l.Parents[i].Name)
			}
			// objects

			if len(*l.Name) > 0 {
				objects = append(objects, *l.Name)
			}
			// TODO face recognition
		}

	}

	imgInfo.persons = nil
	imgInfo.parents = deduplicate(parents)
	imgInfo.objects = objects
	return imgInfo
}

func sendSns(info ImgInfo) (response string, err error) {
	var snsArn = "arn:aws:sns:" + region + ":" + account + ":AddImgTopic"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sns.New(sess)

	result, err := svc.Publish(&sns.PublishInput{
		Subject:  aws.String("SUBJECT - DUMMY"),
		Message:  aws.String("MESSAGE - DUMMY"),
		TopicArn: aws.String(snsArn),
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return *result.MessageId, nil
}

func main() {
	lambda.Start(HandleRequest)
}

type ImgInfo struct {
	persons []string
	parents []string
	objects []string
}
