package environment_variable

import (
	"os"
	"strings"
)

var _ = loadEnv()
var EnvironmentVariable = make(map[string]string)

func loadEnv() bool {
	for _, env := range os.Environ() {
		keyValue := strings.SplitN(env, "=", 2)
		EnvironmentVariable[keyValue[0]] = keyValue[1]
	}

	return true
}

func ReloadEnv() {
	loadEnv()
}
