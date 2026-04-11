package kics

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKicsService(t *testing.T) {
	service := GetKicsService(loader.Config{Path: "./test", Kics: &loader.Kics{}})
	_, ok := service.(interfaces.ScanServiceImpl)

	assert.NotNil(t, service)
	assert.True(t, ok)
}
