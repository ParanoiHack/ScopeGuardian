package sync

import (
	"errors"
	"fmt"
	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/models"
	"ScopeGuardian/logger"
	"time"
)

// sleepFunc is the function used to pause between polling attempts. It is a
// package-level variable so that tests can override it to avoid real sleeps.
var sleepFunc = time.Sleep

// pollInitialDelay is how long to wait before the first polling attempt. It gives
// DefectDojo time to start processing the imported scan before we begin querying
// for findings. All scanners use close_old_findings_product_scope=true, which means
// DD closes old findings synchronously on import but creates new ones asynchronously,
// so without this delay the very first poll could observe a stable zero that is not
// yet the final state.
var pollInitialDelay = 3 * time.Second

// pollInterval is the time to wait between consecutive polling attempts.
var pollInterval = 2 * time.Second

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
// GetActiveFindings is designed to be called after SyncResults; the engagement will already
// exist at that point.
//
// Because DefectDojo processes imported scans asynchronously (old findings are closed
// synchronously on import, new findings are created afterwards), this function polls
// GetFindings at regular intervals until the active finding count stabilises between two
// consecutive reads, ensuring the caller always sees the final post-import state.
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
			return pollFindings(ddService, engagement.Id)
		}
	}

	logger.Info(fmt.Sprintf(logInfoNoEngagementFound, branch))
	return nil, errors.New(errEngagementNotFound)
}

// pollFindings waits for DefectDojo to finish processing an imported scan and then returns
// the active findings for the given engagement. It first waits pollInitialDelay to allow
// DD to begin creating new findings, then polls at pollInterval intervals until the count
// is identical in two consecutive reads (stable) or pollMaxRetries attempts are exhausted.
// An API error on any attempt is returned immediately without further retrying.
func pollFindings(ddService defectdojo.DefectDojoService, engagementId int) ([]defectdojo.Finding, error) {
	logger.Info(fmt.Sprintf(logInfoPollingFindings, engagementId))

	sleepFunc(pollInitialDelay)

	prevCount := -1
	var lastFindings []defectdojo.Finding

	for i := 0; i < pollMaxRetries; i++ {
		if i > 0 {
			sleepFunc(pollInterval)
		}

		findings, err := ddService.GetFindings(engagementId, 0, 100, []defectdojo.Finding{})
		if err != nil {
			logger.Error(fmt.Sprintf(logErrorGetFindings, engagementId))
			return nil, errors.New(errGetFindings)
		}

		if len(findings) == prevCount {
			return findings, nil
		}

		prevCount = len(findings)
		lastFindings = findings
	}

	return lastFindings, nil
}

// FilterByActiveFindings returns only the local findings that have a matching active finding
// in DefectDojo, using a content hash as the match key. The hash is computed by
// models.ComputeFindingHash from fields that are preserved without transformation when a
// scan result is imported into DefectDojo and read back via its API
// (Severity, FilePath/SinkFile, Line/SinkLine, Mitigation/Recommendation). This makes the
// filter completely independent of how DefectDojo formats finding titles (e.g. KICS prefixes
// the rule category to the name), what engine produced the finding, and whether a VulnId
// is present.
//
// Two hashes are computed per DD finding:
//   - hash(""|severity|filePath|line|mitigation): matches KICS and Opengrep local findings,
//     which pass an empty vulnId to ComputeFindingHash.
//   - hash(title|severity|filePath|line|mitigation): matches Grype local findings, where
//     the VulnId (CVE/GHSA) equals the DD finding title.
//
// This filtering respects suppressions applied in DefectDojo: any finding marked as false
// positive or accepted risk will be absent from the active set and therefore dropped locally.
func FilterByActiveFindings(local []models.Finding, active []defectdojo.Finding) []models.Finding {
	activeSet := make(map[string]struct{}, len(active)*2)
	for _, f := range active {
		// Hash without vuln-id covers KICS and Opengrep.
		activeSet[models.ComputeFindingHash("", f.Severity, f.FilePath, f.Line, f.Mitigation)] = struct{}{}
		// Hash with title as vuln-id covers Grype (DD title == CVE/GHSA id).
		activeSet[models.ComputeFindingHash(f.Title, f.Severity, f.FilePath, f.Line, f.Mitigation)] = struct{}{}
	}

	filtered := make([]models.Finding, 0, len(local))
	for _, f := range local {
		if _, ok := activeSet[f.Hash]; ok {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

