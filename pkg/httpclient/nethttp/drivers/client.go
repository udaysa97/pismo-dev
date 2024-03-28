package drivers

import (
	"bytes"
	"fmt"
	"net/http"
	"pismo-dev/pkg/httpclient/types"
	"pismo-dev/pkg/logger"
	"time"
)

type NetHttpWrapper struct {
	client *http.Client
}

func NewNetHttpClient(httpClient *http.Client) *NetHttpWrapper {
	return &NetHttpWrapper{
		client: httpClient,
	}
}

func validate(options types.RequestOptions) error {
	if len(options.Url) == 0 {
		return fmt.Errorf("invalid `Url`")
	}
	if options.RetryAttempt < 0 {
		return fmt.Errorf("invalid `RetryAttempt` url")
	}
	if options.RetryFixedInternal < 0 {
		return fmt.Errorf("invalid `RetryFixedInternal` url")
	}
	return nil
}

func (rq *NetHttpWrapper) Get(options types.RequestOptions) (*http.Response, error) {
	if err := validate(options); err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodGet, options.Url, nil)

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	params := req.URL.Query()
	for k, v := range options.QueryParams {
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= options.RetryAttempt; attempt++ {
		logger.Info(fmt.Sprintf("[GET] %s", options.Url), map[string]any{"attempt": attempt})
		resp, err = rq.client.Do(req)
		if err == nil && resp.StatusCode != 503 {
			break
		}
		logger.Error(fmt.Sprintf("[GET] %s", options.Url), map[string]any{"error": err, "retryFixedInternal": options.RetryFixedInternal})
		time.Sleep(options.RetryFixedInternal)
	}
	return resp, err
}

func (rq *NetHttpWrapper) Post(options types.RequestOptions) (*http.Response, error) {
	if err := validate(options); err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(http.MethodPost, options.Url, bytes.NewBuffer(options.Body))

	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	params := req.URL.Query()
	for k, v := range options.QueryParams {
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()

	var resp *http.Response
	var err error
	for attempt := 0; attempt <= options.RetryAttempt; attempt++ {
		logger.Info(fmt.Sprintf("[POST] %s", options.Url), map[string]any{"attempt": attempt})
		resp, err = rq.client.Do(req)
		if err == nil && resp.StatusCode != 503 {
			break
		}
		logger.Error(fmt.Sprintf("[POST] %s", options.Url), map[string]any{"error": err, "retryFixedInternal": options.RetryFixedInternal})
		time.Sleep(options.RetryFixedInternal)
	}
	return resp, err
}
