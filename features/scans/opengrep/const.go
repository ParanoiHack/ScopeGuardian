package opengrep

const (
	binaryPath   = "/opt/opengrep/bin/opengrep"
	dirPath      = "/opt/opengrep"
	outputFolder = "results"
	scannerType  = "SAST"
)

const (
	severityThreshold    = "Info"
	groupByProperty      = "finding_title"
	findingGroupProperty = true
	findingTagProperty   = true
	SASTEngineTag        = "SAST"
	scanType             = "Semgrep JSON Report"
	closeOldFinding      = true
)

const (
	outputNameParameter = "opengrep-result.json"
)

const (
	jsonOutputArgument          = "--json-output="
	ossOnlyArgument             = "--oss-only"
	quietArgument               = "-q"
	skipUnknownExtArgument      = "--skip-unknown-extensions"
	excludeArgument             = "--exclude"
	excludeRuleArgument         = "--exclude-rule"
)

const (
	logErrorDirectoryNotFound = "Cannot find directory [%s]"
	logErrorFileNotFound      = "Cannot find file [%s]"
	logErrorParseResults      = "Cannot parse opengrep results"
	logErrorImportScan        = "Cannot import scan into engagement [%d]"
	logErrorGetTests          = "Cannot get tests for engagement [%d]"
)

const (
	logInfoCommandLine = "Command line invoked [%s]"
)

const (
	errDirectoryNotFound = "directory not found"
)
