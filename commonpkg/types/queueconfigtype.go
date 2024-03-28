package commontypes

import (
	queuewrapper "pismo-dev/pkg/queue"
)

type QueueProperties struct {
	AWSRegion                 string
	Env                       string
	SQSStandardQueueURL       string
	SQSDelaySeconds           int
	SQSMessageRetentionPeriod string
	SQSVisibilityTimeout      int
	SQSWaitTime               int
	SQSMaxMessage             int
	SQSName                   string
}

type QueueWrapperConfig struct {
	WrapperObj queuewrapper.QueueWrapper
	QueueProps QueueProperties
}
