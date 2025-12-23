package exec

import (
	"os"
	"os/exec"
)

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
