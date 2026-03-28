package exec

import (
	"os"
	"os/exec"
)

// Wrap executes the binary at binaryPath with the given args inside dirPath.
// stdout and stderr of the child process are forwarded to the current process's
// stdout and stderr respectively. It returns true on success or false and the
// underlying error if the command exits with a non-zero status.
func Wrap(binaryPath string, dirPath string, args []string) (bool, error) {
	cmd := exec.Command(binaryPath, args...)

	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return false, err
	}

	return true, nil
}
