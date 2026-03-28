package kics

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

// GetKicsService constructs and returns a ScanServiceImpl for the KICS scanner
// using the provided loader configuration.
func GetKicsService(config loader.Kics) interfaces.ScanServiceImpl {
	return newKicsService(config)
}
