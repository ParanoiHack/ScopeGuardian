package engine

const (
	logInfoKicsRegister    = "Kics enabled. Registring scanner for execution"
	logInfoScannerStarting = "Starting %s scanning engine"
	logInfoScannerSuccess  = "%s scanner succeeded"
	logInfoFindingsLoaded  = "%s findings loaded"
	logInfoSyncResult      = "Syncing %s results to DefectDojo"
)

const (
	logErrorScannerFailed              = "%s scanner failed"
	logErrorLoadFinding                = "Cannot load finding for %s scanner"
	logErrorRegisterScanner            = "Cannot register scanner %s"
	logErrorRetrieveEngagementId       = "Cannot retrieve engagement ID for project [%s] branch [%s]"
	logErrorRetrieveDefectDojoFindings = "Cannot retrieve findings from DefectDojo"
)

const (
	kicsScannerName = "Kics (IACST)"
)
