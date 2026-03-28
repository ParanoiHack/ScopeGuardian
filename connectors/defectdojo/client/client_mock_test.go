package client

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMockClient_MockDo_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	expected := &http.Response{StatusCode: http.StatusOK}

	mock.EXPECT().Do(req).Return(expected, nil)

	resp, err := mock.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMockClient_MockDo_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	mock.EXPECT().Do(req).Return(nil, errors.New("transport error"))

	resp, err := mock.Do(req)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func TestMockClient_MockGet_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	expectedBody := []byte(`{"key":"value"}`)
	mock.EXPECT().Get("https://example.com", http.Header{}).Return(expectedBody, http.StatusOK)

	body, status := mock.Get("https://example.com", http.Header{})
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, expectedBody, body)
}

func TestMockClient_MockGet_ReturnsNilBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	mock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, -1)

	body, status := mock.Get("https://example.com", http.Header{})
	assert.Equal(t, -1, status)
	assert.Nil(t, body)
}

func TestMockClient_MockGetHeaders(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	expectedHeaders := http.Header{AuthorizationKey: {"Token test-token"}}
	mock.EXPECT().GetHeaders("test-token").Return(expectedHeaders)

	headers := mock.GetHeaders("test-token")
	assert.Equal(t, expectedHeaders, headers)
}

func TestMockClient_MockPost_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	expectedBody := []byte(`{"id":1}`)
	mock.EXPECT().Post("https://example.com", gomock.Any(), gomock.Any()).Return(expectedBody, http.StatusCreated)

	body, status := mock.Post("https://example.com", []byte(`{}`), http.Header{})
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, expectedBody, body)
}

func TestMockClient_MockPost_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	mock.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, -1)

	body, status := mock.Post("https://example.com", []byte(`{}`), http.Header{})
	assert.Equal(t, -1, status)
	assert.Nil(t, body)
}

func TestMockClient_MockPut_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	expectedBody := []byte(`{"updated":true}`)
	mock.EXPECT().Put("https://example.com", gomock.Any(), gomock.Any()).Return(expectedBody, http.StatusOK)

	body, status := mock.Put("https://example.com", []byte(`{}`), http.Header{})
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, expectedBody, body)
}

func TestMockClient_MockPut_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockClient(ctrl)

	mock.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, -1)

	body, status := mock.Put("https://example.com", []byte(`{}`), http.Header{})
	assert.Equal(t, -1, status)
	assert.Nil(t, body)
}
