package syft

import (
	"errors"
	"fmt"
	"io"
	"os"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/interfaces"
	"ScopeGuardian/domains/models"
	environment_variable "ScopeGuardian/environnement_variable"
	"ScopeGuardian/exec"
	"ScopeGuardian/logger"
	"strings"
)

// execRunner is the function signature used to invoke an external binary.
// It matches exec.Wrap so that tests can substitute a lightweight mock.
type execRunner func(binaryPath string, dirPath string, args []string, stdout io.Writer, stderr io.Writer, extraEnv ...string) (bool, error)

// SyftServiceImpl implements ScanServiceImpl for the Syft SBOM generator.
type SyftServiceImpl struct {
	path                 string
	transitiveLibraries  bool
	excludeTestLibraries bool
	proxyEnv             []string
	runner               execRunner
}

// newSyftService builds a SyftServiceImpl from the scan path, resolving it
// relative to the SCAN_DIR environment variable. transitiveLibraries controls
// whether Syft resolves transitive Java dependencies from Maven Central.
// excludeTestLibraries controls whether test source directories are excluded
// from the Syft filesystem scan (e.g. **/src/test/**).
// proxyEnv is an optional list of "KEY=VALUE" proxy environment variable entries
// (see loader.Proxy.ToEnv) forwarded to the Syft process.
func newSyftService(path string, transitiveLibraries bool, excludeTestLibraries bool, proxyEnv []string) interfaces.ScanServiceImpl {
	return &SyftServiceImpl{
		path:                 fmt.Sprintf("%s/%s", environment_variable.EnvironmentVariable["SCAN_DIR"], path),
		transitiveLibraries:  transitiveLibraries,
		excludeTestLibraries: excludeTestLibraries,
		proxyEnv:             proxyEnv,
		runner:               exec.Wrap,
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

	if s.transitiveLibraries {
		logger.Info(logInfoTransitiveLibraries)
	}

	if s.excludeTestLibraries {
		args = append(args, excludeArgument, testSourcePattern)
		logger.Info(logInfoExcludeTestLibraries)
	}

	logger.Info(fmt.Sprintf(logInfoCommandLine, strings.Join(args, " ")))

	transitiveValue := fmt.Sprintf("%v", s.transitiveLibraries)
	extraEnv := []string{
		fmt.Sprintf("%s=%s", envJavaUseNetwork, transitiveValue),
		fmt.Sprintf("%s=%s", envJavaResolveTransitiveDependencies, transitiveValue),
	}
	extraEnv = append(extraEnv, s.proxyEnv...)
	return s.runner(binaryPath, dirPath, args, os.Stdout, os.Stderr, extraEnv...)
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
