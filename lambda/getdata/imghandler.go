package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/sns"
)

const maxlabels = 100
const minConfidence float64 = 75.000000

var region string
var account string

func init() {
	region = os.Getenv("Region")
	account = os.Getenv("Account")
	fmt.Println(":init: region: ", region)
	if region == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		region = os.Getenv("Region")
		account = os.Getenv("Account")
	}
}

func HandleRequest(ctx context.Context, ev MyEvent) (response string, err error) {

	fmt.Println(":context: ", ctx)
	fmt.Println(":event: ", ev)

	var bucket string = ev.Records[0].S3.Bucket.Name
	var imgName string = ev.Records[0].S3.Object.Key
	var imgpath string = "https://" + bucket + ".s3.amazonaws.com/" + imgName

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

	var info = getImgInfo(resultJson, imgpath)
	snsid, err := sendSns(info)
	return snsid, err
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

func getImgInfo(res *rekognition.DetectLabelsOutput, imgpath string) (imgInfo ImgInfo) {
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

	imgInfo.Persons = nil
	imgInfo.Parents = deduplicate(parents)
	imgInfo.Objects = objects
	imgInfo.ImgPath = imgpath
	return imgInfo
}

func sendSns(info ImgInfo) (snsid string, err error) {
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
	fmt.Println(":sendSns: ", *result.MessageId)
	return *result.MessageId, nil
}

func main() {
	lambda.Start(HandleRequest)
}

type ImgInfo struct {
	Persons []string `json:"persons"`
	Parents []string `json:"title"`
	Objects []string `json:"author"`
	ImgPath string   `json:"imgpath"`
}

type MyEvent struct {
	Records []struct {
		EventVersion string    `json:"eventVersion"`
		EventSource  string    `json:"eventSource"`
		AwsRegion    string    `json:"awsRegion"`
		EventTime    time.Time `json:"eventTime"`
		EventName    string    `json:"eventName"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			XAmzRequestID string `json:"x-amz-request-id"`
			XAmzID2       string `json:"x-amz-id-2"`
		} `json:"responseElements"`
		S3 struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationID string `json:"configurationId"`
			Bucket          struct {
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
				Arn string `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int    `json:"size"`
				ETag      string `json:"eTag"`
				VersionID string `json:"versionId"`
				Sequencer string `json:"sequencer"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}
