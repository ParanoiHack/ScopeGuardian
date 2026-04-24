package opengrep

import (
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/loader"
)

// GetOpenGrepService constructs and returns a ScanServiceImpl for the OpenGrep scanner
// using the provided loader configuration.
func GetOpenGrepService(config loader.Config) interfaces.ScanServiceImpl {
	return newOpenGrepService(config.Path, *config.Opengrep, config.Proxy.ToEnv())
}
