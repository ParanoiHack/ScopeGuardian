package syft

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSyftService(t *testing.T) {
	service := GetSyftService(loader.Grype{})
	_, ok := service.(interfaces.ScanServiceImpl)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
