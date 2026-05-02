package syft

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided loader configuration. When a Grype configuration is present,
// its TransitiveLibraries flag, SyftExclude patterns, and SyftDepth are forwarded
// to the Syft service to control whether transitive Java dependencies are resolved
// from Maven Central, which filesystem paths are excluded during SBOM generation,
// and how many parent POM levels are recursively resolved.
func GetSyftService(config loader.Config) interfaces.ScanServiceImpl {
	transitiveLibraries := false
	var exclude []string
	depth := 1
	if config.Grype != nil {
		transitiveLibraries = config.Grype.TransitiveLibraries
		exclude = config.Grype.SyftExclude
		if config.Grype.SyftDepth != 0 {
			depth = config.Grype.SyftDepth
		}
	}
	return newSyftService(config.Path, transitiveLibraries, exclude, depth, config.Proxy.ToEnv())
}
