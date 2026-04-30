package syft

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided loader configuration. When a Grype configuration is present,
// its TransitiveLibraries and ExcludeTestLibraries flags are forwarded to the Syft
// service to control whether transitive Java dependencies are resolved from Maven
// Central and whether test source directories are excluded during SBOM generation.
func GetSyftService(config loader.Config) interfaces.ScanServiceImpl {
	transitiveLibraries := false
	excludeTestLibraries := false
	if config.Grype != nil {
		transitiveLibraries = config.Grype.TransitiveLibraries
		excludeTestLibraries = config.Grype.ExcludeTestLibraries
	}
	return newSyftService(config.Path, transitiveLibraries, excludeTestLibraries, config.Proxy.ToEnv())
}
