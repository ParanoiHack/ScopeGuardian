package engine

const (
	logInfoKicsRegister      = "Kics enabled. Registring scanner for execution"
	logInfoSyftRegister      = "Grype enabled. Registering Syft SBOM scanner for execution"
	logInfoGrypeRegister     = "Grype enabled. Registering Grype vulnerability scanner for execution"
	logInfoOpenGrepRegister  = "OpenGrep enabled. Registering OpenGrep SAST scanner for execution"
	logInfoScannerStarting   = "Starting %s scanning engine"
	logInfoScannerSuccess    = "%s scanner succeeded"
	logInfoFindingsLoaded    = "%s findings loaded"
	logInfoSyncResult        = "Syncing %s results to DefectDojo"
	logInfoSyncResultSuccess = "Successfully synced %s results to DefectDojo"
)

const (
	logErrorScannerFailed        = "%s scanner failed"
	logErrorLoadFinding          = "Cannot load finding for %s scanner"
	logErrorRegisterScanner      = "Cannot register scanner %s"
	logErrorRetrieveEngagementId = "Cannot retrieve engagement ID for project [%s] branch [%s]"
	logErrorSkippingScanner      = "Skipping %s scanner because prerequisite %s failed"
	logErrorSyncResult           = "Failed to sync %s results to DefectDojo"
)

const (
	kicsScannerName      = "Kics (IACST)"
	syftScannerName      = "Syft (SBOM)"
	grypeScannerName     = "Grype (SCA)"
	opengrepScannerName  = "OpenGrep (SAST)"
)
