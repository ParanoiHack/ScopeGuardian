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

	t.Run("Should use excludeTestLibraries from grype config when false", func(t *testing.T) {
		service := GetSyftService(loader.Config{Grype: &loader.Grype{ExcludeTestLibraries: false}})
		svc, ok := service.(*SyftServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
		assert.False(t, svc.excludeTestLibraries)
	})

	t.Run("Should use excludeTestLibraries from grype config when true", func(t *testing.T) {
		service := GetSyftService(loader.Config{Grype: &loader.Grype{ExcludeTestLibraries: true}})
		svc, ok := service.(*SyftServiceImpl)

		assert.NotNil(t, service)
		assert.True(t, ok)
		assert.True(t, svc.excludeTestLibraries)
	})
}
