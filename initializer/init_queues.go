package initializer

import (
	"context"
	commontypes "pismo-dev/commonpkg/types"
	"pismo-dev/internal/appconfig"
	"pismo-dev/internal/queueutil"
	"pismo-dev/pkg/logger"
	queuewrapper "pismo-dev/pkg/queue"
	qdrivers "pismo-dev/pkg/queue/drivers"
)

func InitSQS() map[string]commontypes.QueueWrapperConfig {

	queues := appconfig.GetQueues()

	var queueWrappers map[string]commontypes.QueueWrapperConfig = make(map[string]commontypes.QueueWrapperConfig)

	for key, queue := range queues {
		ctx := context.TODO()
		AWSRegion := queue.AWSRegion
		Env := queue.Env
		SQSStandardQueueURL := queue.SQSStandardQueueURL
		SQSDelaySeconds := queue.SQSDelaySeconds
		SQSMessageRetentionPeriod := queue.SQSMessageRetentionPeriod
		SQSVisibilityTimeout := queue.SQSVisibilityTimeout
		SQSName := queue.SQSName
		logger.Info("SQS Initiating", map[string]interface{}{"region": AWSRegion})
		sqsClient, queueURL, err := queueutil.NewSQSClient(ctx, AWSRegion, Env, SQSStandardQueueURL, SQSDelaySeconds, SQSMessageRetentionPeriod, SQSVisibilityTimeout, SQSName)
		if err != nil {
			logger.Error("SQS init failure: %s", err.Error())
			panic(err)
		}
		logger.Info("Queue details:", map[string]interface{}{"region": AWSRegion})
		sqsWriter := qdrivers.NewSQSProducer(ctx, sqsClient, queueURL)
		sqsConsumer := qdrivers.NewSQSConsumer(ctx, sqsClient, *queueURL)
		queueWrappers[key] = commontypes.QueueWrapperConfig{WrapperObj: queuewrapper.NewQueueWrapper(sqsConsumer, sqsWriter), QueueProps: queue}

	}
	return queueWrappers

}
