package httpclient

import (
	"net/http"
	"pismo-dev/pkg/httpclient/types"
)

type HttpClient interface {
	Get(options types.RequestOptions) (*http.Response, error)
	Post(options types.RequestOptions) (*http.Response, error)
}
