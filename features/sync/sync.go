package sync

import (
	"errors"
	"fmt"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/models"
	"ScopeGuardian/logger"
	"time"
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

// GetActiveFindings fetches the active (non-suppressed) findings for the given project and
// branch from DefectDojo. Unlike GetEngagementId it never creates an engagement: if no
// matching engagement exists for the branch it returns an error so the caller can fall back
// to displaying all local findings unfiltered. This read-only behaviour is intentional —
// GetActiveFindings is called before SyncResults so it must not mutate DD state.
func GetActiveFindings(ddService defectdojo.DefectDojoService, projectName string, branch string, protectedBranches []string) ([]defectdojo.Finding, error) {
	product, err := ddService.GetProductByName(projectName)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorGetProduct, projectName))
		return nil, errors.New(errGetProduct)
	}

	engagements, err := ddService.GetEngagements(uint(product.Id), 0, 100, []defectdojo.Engagement{})
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorGetEngagements, product.Id))
		return nil, errors.New(errGetEngagements)
	}

	expectedName := fmt.Sprintf("%s-%s", projectName, branch)
	for _, engagement := range engagements {
		if engagement.Name == expectedName {
			findings, err := ddService.GetFindings(engagement.Id, 0, 100, []defectdojo.Finding{})
			if err != nil {
				logger.Error(fmt.Sprintf(logErrorGetFindings, engagement.Id))
				return nil, errors.New(errGetFindings)
			}
			return findings, nil
		}
	}

	logger.Info(fmt.Sprintf(logInfoNoEngagementFound, branch))
	return nil, errors.New(errEngagementNotFound)
}

// FilterByActiveFindings returns only the local findings that have a matching active finding
// in DefectDojo. A local finding matches a DD finding when their titles are equal. File path
// and line number are intentionally excluded from the match key because scanners (e.g. KICS)
// may emit relative paths that DefectDojo normalises differently when storing, which would
// cause false mismatches and incorrectly drop active findings.
// For Grype findings the VulnId (CVE/GHSA) is used as the title key since DD stores the
// vulnerability ID — not the artifact name — as the finding title; for all other scanners
// the Name field is used directly.
// This filtering respects suppressions applied in DefectDojo: any finding marked as false
// positive or accepted risk will be absent from the active set and therefore dropped locally.
func FilterByActiveFindings(local []models.Finding, active []defectdojo.Finding) []models.Finding {
	activeSet := make(map[string]struct{}, len(active))
	for _, f := range active {
		activeSet[f.Title] = struct{}{}
	}

	filtered := make([]models.Finding, 0, len(local))
	for _, f := range local {
		title := f.VulnId
		if title == "" {
			title = f.Name
		}
		if _, ok := activeSet[title]; ok {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

