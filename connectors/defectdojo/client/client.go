package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"scope-guardian/logger"
)

// Client is the interface for making HTTP requests to the DefectDojo API.
type Client interface {
	// Do executes the given HTTP request and returns the response.
	Do(req *http.Request) (*http.Response, error)
	// Post sends a POST request to url with the given body and headers,
	// returning the response body and HTTP status code.
	Post(url string, body []byte, headers http.Header) ([]byte, int)
	// Put sends a PUT request to url with the given body and headers,
	// returning the response body and HTTP status code.
	Put(url string, body []byte, headers http.Header) ([]byte, int)
	// Get sends a GET request to url with the given headers,
	// returning the response body and HTTP status code.
	Get(url string, headers http.Header) ([]byte, int)
	// GetHeaders builds the standard HTTP headers required by the DefectDojo API,
	// including the Authorization token.
	GetHeaders(accessToken string) http.Header
}

type clientImpl struct {
	client *http.Client
}

// NewClient wraps the provided *http.Client and returns a Client implementation.
func NewClient(client *http.Client) Client {
	return &clientImpl{
		client: client,
	}
}

// Do executes the given HTTP request using the underlying http.Client.
func (c *clientImpl) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// Post sends an HTTP POST request to the given url with body and headers.
// It returns the response body bytes and the HTTP status code.
// Returns nil and -1 on any request, transport, or read error.
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

// Put sends an HTTP PUT request to the given url with body and headers.
// It returns the response body bytes and the HTTP status code.
// Returns nil and -1 on any request, transport, or read error.
func (c *clientImpl) Put(url string, body []byte, headers http.Header) ([]byte, int) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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

// Get sends an HTTP GET request to the given url with the provided headers.
// It returns the response body bytes and the HTTP status code.
// Returns nil and -1 on any request, transport, or read error.
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

// GetHeaders returns a pre-populated http.Header containing the Accept,
// Content-Type, and Authorization (Token) headers required by the DefectDojo API.
func (c *clientImpl) GetHeaders(accessToken string) http.Header {
	return http.Header{
		AcceptKey:        {AcceptValue},
		ContentTypeKey:   {ContentTypeValue},
		AuthorizationKey: {fmt.Sprintf("Token %s", accessToken)},
	}
}
