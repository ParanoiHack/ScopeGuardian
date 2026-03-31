package grype

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

// GetGrypeService constructs and returns a ScanServiceImpl for the Grype vulnerability scanner
// using the provided loader configuration.
func GetGrypeService(config loader.Config) interfaces.ScanServiceImpl {
	return newGrypeService(*config.Grype)
}
