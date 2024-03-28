package downstreamerror

type DownstreamError struct {
	Code        string
	Message     string
	Endpoint    string
	ServiceName string
	HttpStatus  int
}

func (d *DownstreamError) Error() string {
	return d.Code + " " + d.Message + " " + d.ServiceName + " " + d.Endpoint
}

func New(code string, httpStatus int, serviceName string, endpoint string, message string) error {
	downstreamErr := DownstreamError{Code: code, Message: message, HttpStatus: httpStatus, ServiceName: serviceName, Endpoint: endpoint}
	return &downstreamErr
}
