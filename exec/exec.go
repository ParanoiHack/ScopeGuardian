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
	return WrapAllowExitCodes(binaryPath, dirPath, args, stdout, stderr, nil, extraEnv...)
}

// WrapAllowExitCodes is like Wrap but treats any exit code listed in
// successCodes as a successful execution in addition to the standard exit 0.
// This is needed for scanners such as Grype (exit 1 = vulnerabilities found)
// and OpenGrep (exit 2 = findings found) that use non-zero exit codes to
// signal normal "findings present" conditions rather than errors.
func WrapAllowExitCodes(binaryPath string, dirPath string, args []string, stdout io.Writer, stderr io.Writer, successCodes []int, extraEnv ...string) (bool, error) {
	cmd := exec.Command(binaryPath, args...)

	cmd.Dir = dirPath
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			for _, allowed := range successCodes {
				if code == allowed {
					return true, nil
				}
			}
		}
		return false, err
	}

	return true, nil
}
