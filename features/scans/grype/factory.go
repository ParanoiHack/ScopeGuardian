package grype

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetGrypeService constructs and returns a ScanServiceImpl for the Grype vulnerability scanner
// using the provided loader configuration.
func GetGrypeService(config loader.Config) interfaces.ScanServiceImpl {
	return newGrypeService(*config.Grype)
}
