package defectdojo

import (
	"ScopeGuardian/connectors/defectdojo/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefectDojoService(t *testing.T) {
	service := GetDefectDojoService(client.NewClient(&http.Client{}), "http://localhost", "accessToken")
	_, ok := service.(DefectDojoService)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
