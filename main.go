package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	//"os"
	//"log"
	"os"
)

func main() {
	sess := session.New()
	svc := sqs.New(sess)
	queueName := "dev-enablement-keith"
	messageText := "dave"

	switch os.Args[1] {
	case "write":
		qUrl, err := createOurQueue(svc, queueName)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		err = addOurMessage(svc, qUrl, messageText)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
	case "read":
		qURL, err := getQueueUrl(svc, queueName)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		MessageRead, err := readMessage(svc, qURL)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		err = deleteQueue(svc, qURL)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}
		fmt.Println(MessageRead)
	default:
		fmt.Println("Please enter read or write as arguments")
	}
}

func addOurMessage(svc sqsiface.SQSAPI, qURL string, messageText string) error {
	var delaySeconds  int64 = 10

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(delaySeconds),
		MessageBody:  aws.String(messageText),
		QueueUrl:     &qURL,
	})

	if err != nil {
		return err
	}
	fmt.Println("Success added: ", messageText)
	return nil
}

func createOurQueue(svc sqsiface.SQSAPI, queueName string) (string, error) {
	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	if err != nil {
		return "", err
	}
	return *result.QueueUrl, nil
}

func readMessage(svc sqsiface.SQSAPI, qURL string) (string, error) {

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &qURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Error", err)
		return "", err
	}

	if len(result.Messages) == 0 {
		//fmt.Println("Received no messages")
		return "No message on queue", nil
	}

	return *result.Messages[0].Body, nil
}

func deleteQueue(svc sqsiface.SQSAPI, qURL string) error {
	_, err := svc.DeleteQueue(&sqs.DeleteQueueInput{
		QueueUrl: aws.String(qURL),
	})
	if err != nil {
		return err
	}
	return nil
}

func getQueueUrl(svc sqsiface.SQSAPI, queueName string) (string, error) {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	return *result.QueueUrl, nil
}
