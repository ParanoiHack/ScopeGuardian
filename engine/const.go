package engine

const (
	logInfoKicsRegister    = "Kics enabled. Registring scanner for execution"
	logInfoSyftRegister    = "Grype enabled. Registering Syft SBOM scanner for execution"
	logInfoGrypeRegister   = "Grype enabled. Registering Grype vulnerability scanner for execution"
	logInfoScannerStarting = "Starting %s scanning engine"
	logInfoScannerSuccess  = "%s scanner succeeded"
	logInfoFindingsLoaded  = "%s findings loaded"
	logInfoSyncResult      = "Syncing %s results to DefectDojo"
)

const (
	logErrorScannerFailed        = "%s scanner failed"
	logErrorLoadFinding          = "Cannot load finding for %s scanner"
	logErrorRegisterScanner      = "Cannot register scanner %s"
	logErrorRetrieveEngagementId = "Cannot retrieve engagement ID for project [%s] branch [%s]"
)

const (
	kicsScannerName  = "Kics (IACST)"
	syftScannerName  = "Syft (SBOM)"
	grypeScannerName = "Grype (SCA)"
)
