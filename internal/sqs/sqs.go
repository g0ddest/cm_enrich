package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewSession(endpoint, region string) *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	return sess
}

func ReceiveMessages(sqsSvc *sqs.SQS, queueURL string) ([]*sqs.Message, error) {
	result, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
	})

	if err != nil {
		return nil, err
	}

	return result.Messages, nil
}

func DeleteMessage(sqsSvc *sqs.SQS, queueURL string, msg *sqs.Message) error {
	_, err := sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}
