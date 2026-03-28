package sync

import (
	"errors"
	"fmt"
	"scope-guardian/connectors/defectdojo"
	"scope-guardian/logger"
	"time"
)

const (
	logInfoEngagementFound       = "Found existing engagement [%s] with id [%d]"
	logInfoEngagementEndDatePast = "Engagement [%d] end date is in the past, updating"
	logInfoEngagementNotFound    = "No engagement found for branch [%s], creating new one"
	logErrorGetProduct           = "Cannot retrieve product for project [%s]"
	logErrorGetEngagements       = "Cannot retrieve engagements for product [%d]"
	logErrorUpdateEndDate        = "Cannot update engagement end date [%d]"
	logErrorCreateEngagement     = "Cannot create engagement for branch [%s]"
)

const (
	errGetProduct       = "cannot retrieve product"
	errGetEngagements   = "cannot retrieve engagements"
	errUpdateEndDate    = "cannot update engagement end date"
	errCreateEngagement = "cannot create engagement"
)

// GetEngagementId retrieves the engagement ID to use for syncing scan results.
// It looks up the product by name, then searches for an existing engagement matching
// the expected name for the given branch. If the engagement exists but its end date is
// in the past, the end date is updated. If no matching engagement exists, a new one is created.
// protectedBranches lists branches whose engagements receive a one-year end date; all others
// receive one week.
func GetEngagementId(ddService defectdojo.DefectDojoService, projectName string, branch string, protectedBranches []string) (int, error) {
	product, err := ddService.GetProductByName(projectName)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorGetProduct, projectName))
		return 0, errors.New(errGetProduct)
	}

	engagements, err := ddService.GetEngagements(uint(product.Id), 0, 100, []defectdojo.Engagement{})
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorGetEngagements, product.Id))
		return 0, errors.New(errGetEngagements)
	}

	isProtected := isProtectedBranch(branch, protectedBranches)
	expectedName := fmt.Sprintf("%s-%s", projectName, branch)

	for _, engagement := range engagements {
		if engagement.Name == expectedName {
			logger.Info(fmt.Sprintf(logInfoEngagementFound, engagement.Name, engagement.Id))

			endDate, parseErr := time.Parse(defectdojo.DateFormat, engagement.TargetEnd)
			if parseErr == nil && endDate.Before(time.Now()) {
				logger.Info(fmt.Sprintf(logInfoEngagementEndDatePast, engagement.Id))
				_, err = ddService.UpdateEngagementEndDate(engagement.Id, product.Id, isProtected)
				if err != nil {
					logger.Error(fmt.Sprintf(logErrorUpdateEndDate, engagement.Id))
					return 0, errors.New(errUpdateEndDate)
				}
			}

			return engagement.Id, nil
		}
	}

	logger.Info(fmt.Sprintf(logInfoEngagementNotFound, branch))

	engagementId, err := ddService.CreateEngagement(projectName, branch, product.Id, isProtected)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorCreateEngagement, branch))
		return 0, errors.New(errCreateEngagement)
	}

	return engagementId, nil
}

// isProtectedBranch returns true if branch appears in the protectedBranches list.
func isProtectedBranch(branch string, protectedBranches []string) bool {
	for _, b := range protectedBranches {
		if b == branch {
			return true
		}
	}
	return false
}

