package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func newMockHTTPClient(doFunc func(req *http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: &mockRoundTripper{doFunc: doFunc},
	}
}

type mockRoundTripper struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestClient_Post(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		body := io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`)))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil
	})

	client := NewClient(mockClient)

	body, statusCode := client.Post("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, `{"key": "value"}`, string(body))
}

func TestClient_Post_Error(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Post("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClientUtils_Get(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		body := io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`)))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil
	})

	client := NewClient(mockClient)

	body, statusCode := client.Get("https://example.com", http.Header{
		"Authorization": {"Bearer token"},
	})

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, `{"key": "value"}`, string(body))
}

func TestClient_Get_Error(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Get("https://example.com", http.Header{
		"Authorization": {"Bearer token"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_GetHeaders(t *testing.T) {
	client := NewClient(nil)

	headers := client.GetHeaders("token")

	assert.Equal(t, "application/json; version=1.0", headers.Get("Accept"))
	assert.Equal(t, "Token token", headers.Get("Authorization"))
}

func TestClient_Post_RequestError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("request creation error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Post("://invalid-url", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Post_NetworkError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Post("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Post_ReadBodyError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(&errorReader{}), // Simule une erreur de lecture du corps
		}, nil
	})

	client := NewClient(mockClient)

	body, statusCode := client.Post("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Get_RequestError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("request creation error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Get("://invalid-url", http.Header{
		"Authorization": {"Token token"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Get_NetworkError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Get("https://example.com", http.Header{
		"Authorization": {"Token token"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Get_ReadBodyError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(&errorReader{}),
		}, nil
	})

	client := NewClient(mockClient)

	body, statusCode := client.Get("https://example.com", http.Header{
		"Authorization": {"Token token"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *errorReader) Close() error {
	return nil
}

func TestClient_Put(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		body := io.NopCloser(bytes.NewReader([]byte(`{"key": "value"}`)))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil
	})

	client := NewClient(mockClient)

	body, statusCode := client.Put("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, `{"key": "value"}`, string(body))
}

func TestClient_Put_Error(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Put("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Put_RequestError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("request creation error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Put("://invalid-url", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Put_NetworkError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	client := NewClient(mockClient)

	body, statusCode := client.Put("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}

func TestClient_Put_ReadBodyError(t *testing.T) {
	mockClient := newMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(&errorReader{}), // Simule une erreur de lecture du corps
		}, nil
	})

	client := NewClient(mockClient)

	body, statusCode := client.Put("https://example.com", []byte(`{"data":"test"}`), http.Header{
		"Content-Type": {"application/json"},
	})

	assert.Equal(t, -1, statusCode)
	assert.Nil(t, body)
}
