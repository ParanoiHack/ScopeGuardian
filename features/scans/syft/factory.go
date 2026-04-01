package syft

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided loader configuration. When a Grype configuration is present,
// its TransitiveLibraries flag is forwarded to the Syft service to control whether
// transitive Java dependencies are resolved from Maven Central during SBOM generation.
func GetSyftService(config loader.Config) interfaces.ScanServiceImpl {
	transitiveLibraries := false
	if config.Grype != nil {
		transitiveLibraries = config.Grype.TransitiveLibraries
	}
	return newSyftService(config.Path, transitiveLibraries)
}
