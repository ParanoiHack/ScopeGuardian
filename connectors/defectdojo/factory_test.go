package defectdojo

import (
	"net/http"
	"scope-guardian/connectors/defectdojo/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefectDojoService(t *testing.T) {
	service := GetDefectDojoService(client.NewClient(&http.Client{}), "http://localhost", "accessToken")
	_, ok := service.(DefectDojoService)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
