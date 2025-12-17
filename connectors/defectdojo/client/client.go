package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"scope-guardian/logger"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url string, body []byte, headers http.Header) ([]byte, int)
	Get(url string, headers http.Header) ([]byte, int)
	GetHeaders(accessToken string) http.Header
}

type clientImpl struct {
	client *http.Client
}

func NewClient(client *http.Client) Client {
	return &clientImpl{
		client: client,
	}
}

func (c *clientImpl) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *clientImpl) Post(url string, body []byte, headers http.Header) ([]byte, int) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		logger.Error(err.Error())
		return nil, -1
	}

	for key, value := range headers {
		req.Header[key] = value
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil, -1
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, -1
	}

	return body, resp.StatusCode
}

func (c *clientImpl) Get(url string, headers http.Header) ([]byte, int) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, -1
	}

	for key, value := range headers {
		req.Header[key] = value
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil, -1
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, -1
	}

	return body, resp.StatusCode
}

func (c *clientImpl) GetHeaders(accessToken string) http.Header {
	return http.Header{
		AcceptKey:        {AcceptValue},
		AuthorizationKey: {fmt.Sprintf("Token %s", accessToken)},
	}
}
