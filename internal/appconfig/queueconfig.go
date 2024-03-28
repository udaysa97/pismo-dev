package appconfig

import (
	commontypes "pismo-dev/commonpkg/types"
	"pismo-dev/constants"
)

func GetQueues() map[string]commontypes.QueueProperties {
	var queues = map[string]commontypes.QueueProperties{
		constants.RECONQUEUE: {
			AWSRegion:                 AWS_REGION,
			Env:                       SQS_ENV,
			SQSStandardQueueURL:       SQS_STANDARD_QUEUE_URL,
			SQSMaxMessage:             1, //TODO: Make configurable from env
			SQSDelaySeconds:           SQS_DELAY_SECONDS,
			SQSMessageRetentionPeriod: SQS_MESSAGGE_RETENTION_PERIOD,
			SQSVisibilityTimeout:      SQS_VISIBILITY_TIMEOUT,
			SQSWaitTime:               3, //TODO: Make configurable from env
			SQSName:                   SQS_QUEUE_NAME,
		},
		constants.POLYGON_NFTPORT_UID: {
			AWSRegion:                 AWS_REGION,
			Env:                       SQS_ENV,
			SQSStandardQueueURL:       SQS_STANDARD_QUEUE_URL_NFT_PORT_POLYGON,
			SQSDelaySeconds:           SQS_DELAY_SECONDS_NFT_PORT_POLYGON,
			SQSMessageRetentionPeriod: SQS_MESSAGGE_RETENTION_PERIOD_NFT_PORT_POLYGON,
			SQSVisibilityTimeout:      SQS_VISIBILITY_TIMEOUT_NFT_PORT_POLYGON,
			SQSWaitTime:               3,
			SQSMaxMessage:             1,
			SQSName:                   SQS_QUEUE_NAME_NFT_PORT_POLYGON,
		},
		constants.POLYGON_CROSSMINT_UID: {
			AWSRegion:                 AWS_REGION,
			Env:                       SQS_ENV,
			SQSStandardQueueURL:       SQS_STANDARD_QUEUE_URL_CROSSMINT_POLYGON,
			SQSDelaySeconds:           SQS_DELAY_SECONDS_CROSSMINT_POLYGON,
			SQSMessageRetentionPeriod: SQS_MESSAGGE_RETENTION_PERIOD_CROSSMINT_POLYGON,
			SQSVisibilityTimeout:      SQS_VISIBILITY_TIMEOUT_CROSSMINT_POLYGON,
			SQSWaitTime:               3,
			SQSMaxMessage:             1,
			SQSName:                   SQS_QUEUE_NAME_CROSSMINT_POLYGON,
		},
	}
	return queues
}
