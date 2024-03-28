package models

type Status string

var (
	StatusSuccess Status = "SUCCESS"
	StatusError   Status = "ERROR"
	StatusFail    Status = "FAIL"
)

type ErrorResponse struct {
	Status  Status `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}
