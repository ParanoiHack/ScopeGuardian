package grype

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGrypeService(t *testing.T) {
	service := GetGrypeService(loader.Config{Grype: &loader.Grype{}})
	_, ok := service.(interfaces.ScanServiceImpl)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
