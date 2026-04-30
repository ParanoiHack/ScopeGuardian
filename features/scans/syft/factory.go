package syft

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided loader configuration. When a Grype configuration is present,
// its TransitiveLibraries flag and SyftExclude patterns are forwarded to the Syft
// service to control whether transitive Java dependencies are resolved from Maven
// Central and which filesystem paths are excluded during SBOM generation.
func GetSyftService(config loader.Config) interfaces.ScanServiceImpl {
	transitiveLibraries := false
	var exclude []string
	if config.Grype != nil {
		transitiveLibraries = config.Grype.TransitiveLibraries
		exclude = config.Grype.SyftExclude
	}
	return newSyftService(config.Path, transitiveLibraries, exclude, config.Proxy.ToEnv())
}
