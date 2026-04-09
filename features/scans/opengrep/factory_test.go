package opengrep

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOpenGrepService(t *testing.T) {
	service := GetOpenGrepService(loader.Config{Path: "./test", Opengrep: &loader.Opengrep{}})
	_, ok := service.(interfaces.ScanServiceImpl)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
