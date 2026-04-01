package exec

import (
	"io"
	"os"
	"os/exec"
)

// Wrap executes the binary at binaryPath with the given args inside dirPath.
// stdout and stderr of the child process are written to the provided io.Writer
// values. Pass os.Stdout / os.Stderr to forward output to the terminal, or
// io.Discard to suppress it. It returns true on success or false and the
// underlying error if the command exits with a non-zero status.
// Optional extraEnv entries (formatted as "KEY=VALUE") are appended to the
// child process environment without affecting the parent process.
func Wrap(binaryPath string, dirPath string, args []string, stdout io.Writer, stderr io.Writer, extraEnv ...string) (bool, error) {
	cmd := exec.Command(binaryPath, args...)

	cmd.Dir = dirPath
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	if err := cmd.Run(); err != nil {
		return false, err
	}

	return true, nil
}
