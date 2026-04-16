package defectdojo

const (
	APIPrefix            = "/api/v2"
	GetProductsPath      = "/products?name_exact="
	GetEngagementsPath   = "/engagements?product=%d&offset=%d&limit=%d"
	CreateEngagementPath = "/engagements/"
	UpdateEngagementPath = "/engagements/%d/"
	ImportScanPath       = "/import-scan/"
	ReimportScanPath     = "/reimport-scan/"
	GetTestsPath         = "/tests/?engagement=%d&scan_type=%s"
	GetFindingsPath      = "/findings/?test__engagement=%d&active=true&offset=%d&limit=%d"
)

const (
	FileValueKey = "file"
	FormValuekey = "form"
	DateFormat   = "2006-01-02"
)

const (
	ScopeGuardianTag      = "SCOPE-GUARDIAN"
	EngagementDescription = "Engagement used for security scans' findings storage affecting the %s branch"
	EngagementStatus      = "In Progress"
	EngagementDefaultLead = 1
	EngagementType        = "CI/CD"
)

const (
	logErrorRetrieveProducts        = "Cannot retrieve defectdojo products"
	logErrorRetrieveEngagements     = "Cannot retrieve engagements for product ID %d"
	logErrorAuthorization           = "Wrong API Key"
	logErrorDecodingToken           = "Cannot unmarshall json [%s]"
	logErrorEncodingStruct          = "Cannot marshall payload"
	logErrorCreateEngagement        = "Cannot create engagement"
	logErrorUpdateEngagementEndDate = "Cannot update the end date of the engagement [%d]"
	logErrorReflection              = "Cannot reflect field [%s]"
	logErrorUnknownType             = "Unknow type [%s]"
	logErrorAddFile                 = "Cannot attach file to request"
	logErrorCreateMultipartRequest  = "Cannot create multipart request"
	logErrorImportScan              = "Cannot import scan via %s (HTTP %d)"
	logErrorReimportScan            = "Cannot reimport scan via %s (HTTP %d)"
	logErrorRetrieveTests           = "Cannot retrieve tests for engagement ID %d"
	logErrorRetrieveFindings        = "Cannot retrieve findings for engagement ID %d"
)

const (
	errRetrieveProducts        = "cannot retrieve defectdojo products"
	errRetrieveEngagements     = "cannot retrieve product's engagement"
	errAuthtorization          = "wrong api key"
	errDuplicateProduct        = "two products with the same name already exist"
	errProductNotExist         = "product does not exist"
	errUnmarshal               = "cannot unmarshal json"
	errCreateEngagement        = "cannot create engagement"
	errUpdateEngagementEndDate = "cannot update engagement end date"
	errWritingFile             = "cannot write file to request"
	errImportScan              = "cannot import scan"
	errReimportScan            = "cannot reimport scan"
	errRetrieveTests           = "cannot retrieve tests"
	errRetrieveFindings        = "cannot retrieve findings"
)
