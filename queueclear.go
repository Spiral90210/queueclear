package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"os"
)

func main() {
	if err := doMain(); err != nil {
		fmt.Printf("unhandled error in program: %s\n", err)
		os.Exit(1)
	}
}

func doMain() error {
	ctx := context.Background()

	flagset := flag.NewFlagSet("queueclear", flag.ExitOnError)
	queueNamePtr := flagset.String("queue", "", "Queue to clear")
	quietPtr := flagset.Bool("quiet", false, "Set to not write queue body to stdout")
	_ = flagset.Parse(os.Args[1:])

	if *queueNamePtr == "" {
		return fmt.Errorf("must provide queue name -queue")
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)

	queueUrl, err := sqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: queueNamePtr,
	})
	if err != nil {
		return err
	}

	var softProcessingLimit int64 = 10000
	var msgCount int64

	for msgCount < softProcessingLimit {
		msgData, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            queueUrl.QueueUrl,
			MaxNumberOfMessages: 10,
			VisibilityTimeout:   60, // we won't need it a minute, but why not
		})
		if err != nil {
			return err
		}
		if len(msgData.Messages) == 0 {
			fmt.Printf("processed %d messages\n", msgCount)
			return nil
		}
		for _, msg := range msgData.Messages {
			_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      queueUrl.QueueUrl,
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				return err
			}
			if !*quietPtr {
				_, _ = fmt.Fprintf(os.Stdout, "%s", *msg.Body)
			}
			msgCount += 1
		}
	}
	fmt.Printf("processed %d messages (limited to %d)\n", msgCount, softProcessingLimit)
	return nil
}
