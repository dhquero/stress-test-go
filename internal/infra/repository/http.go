package repository

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

type HTTPOutput struct {
	StatusCode int
}

type HTTPRepository struct {
	url     string
	timeout uint
}

func NewHTTPRepository(url string, timeout uint) *HTTPRepository {
	return &HTTPRepository{
		url:     url,
		timeout: timeout,
	}
}

func (r *HTTPRepository) Get() (*HTTPOutput, error) {
	client := http.Client{
		Timeout: time.Duration(r.timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(r.url)
	if err != nil {
		return &HTTPOutput{StatusCode: 0}, errors.New("request error")
	}

	defer resp.Body.Close()

	return &HTTPOutput{
		StatusCode: resp.StatusCode,
	}, nil
}
