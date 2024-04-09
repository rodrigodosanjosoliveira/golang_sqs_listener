package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const QueueURL = "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/minha-fila"

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}))

	svc := sqs.New(sess)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-signalCh:
			fmt.Println("Exiting...")
			return
		default:
			receiveParams := sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(QueueURL),
				WaitTimeSeconds:     aws.Int64(20),
				MaxNumberOfMessages: aws.Int64(1),
			}

			result, err := svc.ReceiveMessage(&receiveParams)
			if err != nil {
				fmt.Println("Error", err)
				time.Sleep(1 * time.Second)
				continue
			}

			for _, message := range result.Messages {
				fmt.Printf("Message received: %s \n", *message.Body)

				deleteParams := sqs.DeleteMessageInput{
					QueueUrl:      aws.String(QueueURL),
					ReceiptHandle: message.ReceiptHandle,
				}

				_, err := svc.DeleteMessage(&deleteParams)
				if err != nil {
					fmt.Println("Error", err)
					time.Sleep(1 * time.Second)
					continue
				}
			}
		}
	}

}
