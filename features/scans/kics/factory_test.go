package kics

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKicsService(t *testing.T) {
	service := GetKicsService(loader.Kics{})
	_, ok := service.(interfaces.ScanServiceImpl)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
