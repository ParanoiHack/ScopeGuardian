package opengrep

import (
	"scope-guardian/domains/interfaces"
	"scope-guardian/loader"
)

// GetOpenGrepService constructs and returns a ScanServiceImpl for the OpenGrep scanner
// using the provided loader configuration.
func GetOpenGrepService(config loader.Config) interfaces.ScanServiceImpl {
	return newOpenGrepService(config.Path, *config.Opengrep)
}
