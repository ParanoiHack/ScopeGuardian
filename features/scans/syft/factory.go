package syft

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided loader configuration.
func GetSyftService(config loader.Config) interfaces.ScanServiceImpl {
	return newSyftService(config.Path)
}
