package sync

const (
	logInfoEngagementFound       = "Found existing engagement [%s] with id [%d]"
	logInfoEngagementEndDatePast = "Engagement [%d] end date is in the past, updating"
	logInfoEngagementNotFound    = "No engagement found for branch [%s], creating new one"
	logInfoNoEngagementFound     = "No previous DefectDojo engagement found for branch [%s], skipping filter"
	logErrorGetProduct           = "Cannot retrieve product for project [%s]"
	logErrorGetEngagements       = "Cannot retrieve engagements for product [%d]"
	logErrorUpdateEndDate        = "Cannot update engagement end date [%d]"
	logErrorCreateEngagement     = "Cannot create engagement for branch [%s]"
	logErrorGetFindings          = "Cannot retrieve findings from DefectDojo for engagement [%d]"
)

const (
	errGetProduct          = "cannot retrieve product"
	errGetEngagements      = "cannot retrieve engagements"
	errUpdateEndDate       = "cannot update engagement end date"
	errCreateEngagement    = "cannot create engagement"
	errGetFindings         = "cannot retrieve findings"
	errEngagementNotFound  = "no engagement found for this branch"
)
