package types

import (
	"time"
)

type RequestOptions struct {
	PathParams         map[string]string
	QueryParams        map[string]string
	FormData           map[string]string //TODO: not used in net/http impl
	Headers            map[string]string
	RetryAttempt       int
	RetryFixedInternal time.Duration
	Url                string
	Body               []byte
}
