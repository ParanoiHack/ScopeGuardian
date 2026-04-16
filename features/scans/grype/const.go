package grype

const (
	binaryPath  = "/opt/grype/bin/grype"
	dirPath     = "/opt/grype"
	configPath  = "/opt/grype/config/grype.yaml"
	outputFolder = "results"
	scannerType  = "SCA"
)

const (
	severityThreshold    = "Info"
	groupByProperty      = "finding_title"
	findingGroupProperty = true
	findingTagProperty   = true
	SCAEngineTag         = "SCA"
	scanType             = "Anchore Grype"
	closeOldFinding      = true
)

const (
	sbomScheme             = "sbom:"
	sbomInputNameParameter = "sbom.syft.json"
	outputNameParameter    = "grype-result.json"
	outputFormatParameter  = "json"
)

const (
	ignoreStatesArgument = "--ignore-states"
	outputArgument       = "-o"
	quietArgument        = "-q"
	fileArgument         = "--file"
	configArgument       = "-c"
	excludeArgument      = "--exclude"
)

const (
	logErrorSbomNotFound = "Cannot find sbom file [%s]"
	logErrorFileNotFound = "Cannot find file [%s]"
	logErrorParseResults = "Cannot parse grype results"
	logErrorImportScan   = "Cannot import scan into engagement [%d]"
	logErrorGetTests     = "Cannot get tests for engagement [%d]"
)

const (
	logInfoCommandLine = "Command line invoked [%s]"
)

const (
	errSbomNotFound = "sbom not found"
)

const (
	recommendationUpgrade         = "Upgrade to %s"
	recommendationUpgradeMultiple = "Upgrade to one of: %s"
)
