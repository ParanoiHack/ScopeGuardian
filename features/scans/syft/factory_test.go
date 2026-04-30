package syft

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSyftService(t *testing.T) {
	t.Run("Should return a ScanServiceImpl without grype config", func(t *testing.T) {
		service := GetSyftService(loader.Config{})
		_, ok := service.(interfaces.ScanServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
	})

	t.Run("Should use transitiveLibraries from grype config when false", func(t *testing.T) {
		service := GetSyftService(loader.Config{Grype: &loader.Grype{TransitiveLibraries: false}})
		_, ok := service.(interfaces.ScanServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
	})

	t.Run("Should use transitiveLibraries from grype config when true", func(t *testing.T) {
		service := GetSyftService(loader.Config{Grype: &loader.Grype{TransitiveLibraries: true}})
		_, ok := service.(interfaces.ScanServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
	})

	t.Run("Should use syft_exclude from grype config when empty", func(t *testing.T) {
		service := GetSyftService(loader.Config{Grype: &loader.Grype{SyftExclude: nil}})
		svc, ok := service.(*SyftServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
		assert.Empty(t, svc.exclude)
	})

	t.Run("Should use syft_exclude from grype config when set", func(t *testing.T) {
		patterns := []string{"**/src/test/**", "**/testdata/**"}
		service := GetSyftService(loader.Config{Grype: &loader.Grype{SyftExclude: patterns}})
		svc, ok := service.(*SyftServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
		assert.Equal(t, patterns, svc.exclude)
	})
}
