package syft

const (
	binaryPath = "/opt/syft/bin/syft"
	dirPath    = "/opt/syft"
	configPath = "/opt/syft/config/syft.yaml"
)

const (
	scanArgument   = "scan"
	configArgument = "-c"
	outputArgument = "--output"
	quietArgument  = "-q"
)

const (
	outputFolder                     = "results"
	syftJsonOutputNameParameter      = "sbom.syft.json"
	cyclonedxJsonOutputNameParameter = "sbom.cyclonedx.json"
)

const (
	logInfoCommandLine        = "Command line invoked [%s]"
	logErrorDirectoryNotFound = "Cannot find directory [%s]"
)

const (
	errDirectoryNotFound = "directory not found"
)
