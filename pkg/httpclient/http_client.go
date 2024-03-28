package httpclient

type HttpClientWrapper struct {
	DriverName string
	Driver     HttpClient
}

func NewHttpClientWrapper(driverName string, driver HttpClient) *HttpClientWrapper {
	return &HttpClientWrapper{
		DriverName: driverName,
		Driver:     driver,
	}
}
