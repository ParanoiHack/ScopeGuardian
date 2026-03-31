package exec

import (
	"os"
	"os/exec"
)

// Wrap executes the binary at binaryPath with the given args inside dirPath.
// stdout and stderr of the child process are forwarded to the current process's
// stdout and stderr respectively. It returns true on success or false and the
// underlying error if the command exits with a non-zero status.
// Optional extraEnv entries (formatted as "KEY=VALUE") are appended to the
// child process environment without affecting the parent process.
func Wrap(binaryPath string, dirPath string, args []string, extraEnv ...string) (bool, error) {
	cmd := exec.Command(binaryPath, args...)

	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	if err := cmd.Run(); err != nil {
		return false, err
	}

	return true, nil
}
