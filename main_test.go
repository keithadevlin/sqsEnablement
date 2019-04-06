package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
)

const (
	queueName           = "testQueueName"
	mockedQueueUrl      = "www.qurl.com"
	mockedEmptyQueueUrl = "www.empty.com"
	ourMessage          = "TestString"
)

type mockSqsCreateQueue struct {
	sqsiface.SQSAPI
	QueOutput sqs.CreateQueueOutput
}

func (m mockSqsCreateQueue) CreateQueue(input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	m.QueOutput.QueueUrl = aws.String(mockedQueueUrl)
	return &m.QueOutput, nil
}

func TestCreateOurQueue(t *testing.T) {
	mockSqsClient := &mockSqsCreateQueue{}
	actualQUrl, err := createOurQueue(mockSqsClient, queueName)
	assert.NoError(t, err)
	assert.Equal(t, mockedQueueUrl, actualQUrl)
}

type mockSqsQueueAddMessage struct {
	sqsiface.SQSAPI
	MessOutput sqs.SendMessageOutput
}

func (m mockSqsQueueAddMessage) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return &m.MessOutput, nil
}

func TestAddOurMessage(t *testing.T) {
	mockSqsClient := &mockSqsQueueAddMessage{}
	err := addOurMessage(mockSqsClient, mockedQueueUrl, ourMessage)
	assert.NoError(t, err)
}

type mockSqsQueueUrl struct {
	sqsiface.SQSAPI
	QueUrlOutput sqs.GetQueueUrlOutput
}

func (m mockSqsQueueUrl) GetQueueUrl(input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	m.QueUrlOutput.QueueUrl = aws.String(mockedQueueUrl)
	return &m.QueUrlOutput, nil
}

func TestGetQUrl(t *testing.T) {
	mockSqsClient := &mockSqsQueueUrl{}
	actualQUrl, err := getQueueUrl(mockSqsClient, queueName)
	assert.NoError(t, err)
	assert.Equal(t, mockedQueueUrl, actualQUrl)
}

type mockSqsReadMessage struct {
	sqsiface.SQSAPI
	readMessOut sqs.ReceiveMessageOutput
	readMessIn  sqs.ReceiveMessageInput
}

func (m mockSqsReadMessage) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if *m.readMessIn.QueueUrl == mockedQueueUrl {
		m.readMessOut.Messages = []*sqs.Message{{Body: aws.String(ourMessage)}}
	}
	return &m.readMessOut, nil
}

func TestReadMessageSuccess(t *testing.T) {
	mockSqsClient := &mockSqsReadMessage{readMessIn: sqs.ReceiveMessageInput{QueueUrl: aws.String(mockedQueueUrl)}}
	actualMessage, err := readMessage(mockSqsClient, mockedQueueUrl)
	assert.NoError(t, err)
	assert.Equal(t, ourMessage, actualMessage)
}

func TestReadMessageEmptyQueue(t *testing.T) {
	mockSqsClient := &mockSqsReadMessage{readMessIn: sqs.ReceiveMessageInput{QueueUrl: aws.String(mockedEmptyQueueUrl)}}
	actualMessage, err := readMessage(mockSqsClient, mockedQueueUrl)
	assert.NoError(t, err)
	assert.Equal(t, "No message on queue", actualMessage)
}
