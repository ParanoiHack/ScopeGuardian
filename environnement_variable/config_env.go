package environment_variable

import (
	"os"
	"strings"
)

// EnvironmentVariable is a map of all environment variables loaded at startup.
// Keys and values correspond directly to the process environment.
var _ = loadEnv()
var EnvironmentVariable = make(map[string]string)

// loadEnv populates EnvironmentVariable from the current process environment.
// It is called automatically during package initialization.
func loadEnv() bool {
	for _, env := range os.Environ() {
		keyValue := strings.SplitN(env, "=", 2)
		EnvironmentVariable[keyValue[0]] = keyValue[1]
	}

	return true
}

// ReloadEnv re-reads the process environment and updates EnvironmentVariable.
// Useful when environment variables are set after package initialization.
func ReloadEnv() {
	loadEnv()
}
