package types

import (
	"encoding/json"
	"pismo-dev/constants"
)

type SQSConfigJSON struct {
	SQSDelaySeconds           string `json:"sqs_delay_seconds" validate:"required"`
	SQSMessageRetentionPeriod string `json:"sqs_message_retention_period" validate:"required"`
	SQSVisibilityTimeout      string `json:"sqs_visibility_timeout" validate:"required"`
	SQSStandardQueueURL       string `json:"sqs_standard_queue_url" validate:"required"`
}

func (sqs *SQSConfigJSON) Decode(s string) error {
	//default - only dev
	sqs.SQSDelaySeconds = constants.SQSDelaySecondsDefault
	sqs.SQSMessageRetentionPeriod = constants.SQSMessageRetentionPeriodDefault
	sqs.SQSVisibilityTimeout = constants.SQSVisibilityTimeoutDefault
	return json.Unmarshal([]byte(s), sqs)
}
