package syft

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided loader configuration. When a Grype configuration is present,
// its TransitiveLibraries flag, SyftExclude patterns, and SyftMaxParentRecursiveDepth
// are forwarded to the Syft service to control whether transitive Java dependencies
// are resolved from Maven Central, which filesystem paths are excluded during SBOM
// generation, and how many parent POM levels are recursively resolved.
func GetSyftService(config loader.Config) interfaces.ScanServiceImpl {
	transitiveLibraries := false
	var exclude []string
	maxParentRecursiveDepth := 0
	if config.Grype != nil {
		transitiveLibraries = config.Grype.TransitiveLibraries
		exclude = config.Grype.SyftExclude
		maxParentRecursiveDepth = config.Grype.SyftMaxParentRecursiveDepth
	}
	return newSyftService(config.Path, transitiveLibraries, exclude, maxParentRecursiveDepth, config.Proxy.ToEnv())
}
