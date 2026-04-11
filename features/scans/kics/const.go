package kics

const (
	binaryPath   = "/opt/kics/bin/kics"
	dirPath      = "/opt/kics"
	outputFolder = "results"
	scannerType  = "IACST"
)

const (
	severityThreshold    = "Info"
	groupByProperty      = "finding_title"
	findingGroupProperty = true
	findingTagProperty   = true
	IACSTEngineTag       = "IACST"
	scanType             = "KICS Scan"
	closeOldFinding      = true
)

const (
	scanArgument            = "scan"
	ciArgument              = "--ci"
	librariesPathArgument   = "--libraries-path"
	queriesPathArgument     = "--queries-path"
	outputPathArgument      = "--output-path"
	outputNameArgument      = "--output-name"
	pathArgument            = "--path"
	typeArgument            = "--type"
	ignoreOnExitArgument    = "--ignore-on-exit"
	excludeQueriesArgument  = "--exclude-queries"
)

const (
	librariesPathParameter = "/opt/kics/assets/libraries"
	queriesPathParameter   = "/opt/kics/assets/queries"
	outputNameParameter    = "kics-results.json"
	ignoreOnExitParameter  = "results"
)

const (
	logErrorDirectorNotFound = "Cannot find directory [%s]"
	logErrorFileNotFound     = "Cannot find file [%s]"
	logErrorParseResults     = "Cannot parse kics results"
	logErrorImportScan       = "Cannot import scan into engagement [%d]"
)

const (
	logInfoCommandLine = "Command line invoked [%s]"
)

const (
	errDirectoryNotFound = "directory not found"
)
