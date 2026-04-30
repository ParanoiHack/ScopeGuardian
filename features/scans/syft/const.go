package syft

const (
	binaryPath = "/opt/syft/bin/syft"
	dirPath    = "/opt/syft"
	configPath = "/opt/syft/config/syft.yaml"
)

const (
	scanArgument    = "scan"
	configArgument  = "-c"
	outputArgument  = "--output"
	quietArgument   = "-q"
	excludeArgument = "--exclude"
)

const (
	outputFolder                     = "results"
	syftJsonOutputNameParameter      = "sbom.syft.json"
	cyclonedxJsonOutputNameParameter = "sbom.cyclonedx.json"
)

const (
	logInfoCommandLine         = "Command line invoked [%s]"
	logInfoTransitiveLibraries = "Transitive libraries resolution is enabled. This may significantly increase scan time."
	logErrorDirectoryNotFound  = "Cannot find directory [%s]"
)

const (
	errDirectoryNotFound = "directory not found"
)

const (
	envJavaUseNetwork                    = "SYFT_JAVA_USE_NETWORK"
	envJavaResolveTransitiveDependencies = "SYFT_JAVA_RESOLVE_TRANSITIVE_DEPENDENCIES"
)
