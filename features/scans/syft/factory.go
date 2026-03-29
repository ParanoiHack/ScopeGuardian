package syft

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

// GetSyftService constructs and returns a ScanServiceImpl for the Syft SBOM generator
// using the provided Grype loader configuration.
func GetSyftService(config loader.Grype) interfaces.ScanServiceImpl {
	return newSyftService(config)
}
