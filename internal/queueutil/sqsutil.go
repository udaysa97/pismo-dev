package queueutil

import (
	"context"
	"fmt"
	"pismo-dev/constants"
	"pismo-dev/pkg/logger"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go-v2/aws"
)

func NewSQSClient(ctx context.Context, region string, env string, stdQueueUrl string, delaySeconds int, messageRetentionPeriod string, visibilityTimeout int, SQSName string) (client *sqs.Client, qUrl *string, err error) {

	if env == constants.DEVELOPMENT_ENV {
		return NewSQSClientDev(ctx, region, strconv.Itoa(delaySeconds), messageRetentionPeriod, strconv.Itoa(visibilityTimeout), SQSName)
	} else if cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region)); err != nil {
		return nil, nil, err
	} else {
		awstrace.AppendMiddleware(&cfg, awstrace.WithServiceName(constants.SQSTraceDriverName))
		client := sqs.NewFromConfig(cfg)
		logger.Info(fmt.Sprintf(" SQS Client %v and queueUrl %s", client, stdQueueUrl))
		return client, aws.String(stdQueueUrl), nil
	}
}

func CreateSQSQueue(ctx context.Context, queueName string, delaySeconds string, messageRetentionPeriod string, visibilityTimeout string, client *sqs.Client) (url *string, err error) {

	queueAttributes := map[string]string{
		constants.DelaySeconds:           delaySeconds,
		constants.MessageRetentionPeriod: messageRetentionPeriod,
		constants.VisibilityTimeout:      visibilityTimeout,
	}

	input := &sqs.CreateQueueInput{
		QueueName:  aws.String(queueName),
		Attributes: queueAttributes,
	}

	result, err := client.CreateQueue(ctx, input)
	if err != nil {
		return
	} else {
		url = result.QueueUrl
	}

	return
}

func NewSQSClientDev(ctx context.Context, region string, delaySeconds string, messageRetentionPeriod string, visibilityTimeout string, SQSName string) (client *sqs.Client, qUrl *string, err error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://127.0.0.1:4566",
			SigningRegion: region,
		}, nil

	})

	if cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(customResolver)); err != nil {
		return nil, nil, err
	} else {
		client := sqs.NewFromConfig(cfg)
		input := &sqs.GetQueueUrlInput{QueueName: aws.String(SQSName)}
		if queueUrlDetails, err := client.GetQueueUrl(ctx, input); err == nil {
			qUrl := queueUrlDetails.QueueUrl
			return client, qUrl, nil
		} else if qUrl, err := CreateSQSQueue(ctx, SQSName, delaySeconds, messageRetentionPeriod, visibilityTimeout, client); err != nil {
			return client, nil, err
		} else {
			return client, qUrl, nil
		}

	}
}
