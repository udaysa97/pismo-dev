package client

import (
	"net/http"
	"time"
)

func NewDefaultNetHttpClient() *http.Client {
	return &http.Client{}
}

func NewCustomNetHttpClient() *http.Client {
	transport := http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   100,
		MaxIdleConns:          100,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     false,
	}
	httpClient := &http.Client{
		Transport: &transport,
		Timeout:   15 * time.Second,
	}
	return httpClient
}
