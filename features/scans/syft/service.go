package syft

import (
	"errors"
	"fmt"
	"os"
	"scope-guardian/connectors/defectdojo"
	"scope-guardian/domains/interfaces"
	"scope-guardian/domains/models"
	environment_variable "scope-guardian/environnement_variable"
	"scope-guardian/exec"
	"scope-guardian/logger"
	"strings"
)

// SyftServiceImpl implements ScanServiceImpl for the Syft SBOM generator.
type SyftServiceImpl struct {
	path                string
	transitiveLibraries bool
}

// newSyftService builds a SyftServiceImpl from the scan path, resolving it
// relative to the SCAN_DIR environment variable. transitiveLibraries controls
// whether Syft resolves transitive Java dependencies from Maven Central.
func newSyftService(path string, transitiveLibraries bool) interfaces.ScanServiceImpl {
	return &SyftServiceImpl{
		path:                fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], path),
		transitiveLibraries: transitiveLibraries,
	}
}

// Start validates the scan directory and then invokes the Syft binary to generate
// the SBOM. It returns true on success or false and an error if the directory is
// missing or the Syft process exits with a non-zero status.
func (s *SyftServiceImpl) Start() (bool, error) {
	if _, err := os.Stat(s.path); err != nil {
		logger.Error(fmt.Sprintf(logErrorDirectoryNotFound, s.path))
		return false, errors.New(errDirectoryNotFound)
	}

	args := []string{
		scanArgument, s.path,
		configArgument, configPath,
		outputArgument, fmt.Sprintf("syft-json=%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, syftJsonOutputNameParameter),
		outputArgument, fmt.Sprintf("cyclonedx-json=%s/%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], outputFolder, cyclonedxJsonOutputNameParameter),
		quietArgument,
	}

	logger.Info(fmt.Sprintf(logInfoCommandLine, strings.Join(args, " ")))

	transitiveValue := fmt.Sprintf("%v", s.transitiveLibraries)
	return exec.Wrap(binaryPath, dirPath, args,
		fmt.Sprintf("%s=%s", envJavaUseNetwork, transitiveValue),
		fmt.Sprintf("%s=%s", envJavaResolveTransitiveDependencies, transitiveValue),
	)
}

// LoadFindings is intentionally empty: Syft is used only to produce the SBOM
// and does not contribute findings to the security gate.
func (s *SyftServiceImpl) LoadFindings() ([]models.Finding, error) {
	return nil, nil
}

// Sync is intentionally empty: Syft output is consumed downstream by Grype and
// Dependency-Track rather than being imported directly into DefectDojo.
func (s *SyftServiceImpl) Sync(_ int, _ string, _ defectdojo.DefectDojoService) error {
	return nil
}
