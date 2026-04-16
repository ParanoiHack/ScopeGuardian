package sync

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestMain disables real sleeps for all tests in this package so that polling
// in GetActiveFindings/pollFindings completes instantly.
func TestMain(m *testing.M) {
	sleepFunc = func(_ time.Duration) {}
	os.Exit(m.Run())
}

const (
	testProjectName  = "my-project"
	testBranch       = "main"
	testProductId    = 5
	testEngagementId = 42
)

var noProtectedBranches = []string{}
var protectedBranchList = []string{"main", "master"}

func futureDate() string {
	return time.Now().AddDate(1, 0, 0).Format(defectdojo.DateFormat)
}

func pastDate() string {
	return time.Now().AddDate(-1, 0, 0).Format(defectdojo.DateFormat)
}

func expectedEngagementName() string {
	return fmt.Sprintf("%s-%s", testProjectName, testBranch)
}

func TestGetEngagementId(t *testing.T) {
	t.Run("Should create engagement when none exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)
		mockService.EXPECT().CreateEngagement(testProjectName, testBranch, testProductId, false).Return(testEngagementId, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})

	t.Run("Should create engagement with protected=true when branch is in protected list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)
		mockService.EXPECT().CreateEngagement(testProjectName, testBranch, testProductId, true).Return(testEngagementId, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch, protectedBranchList)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})

	t.Run("Should return existing engagement ID when engagement exists with valid end date", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: futureDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})

	t.Run("Should update end date and return engagement ID when end date is in the past", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: pastDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().UpdateEngagementEndDate(testEngagementId, testProductId, false).Return(true, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})

	t.Run("Should update end date with protected=true when branch is protected and end date is past", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: pastDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().UpdateEngagementEndDate(testEngagementId, testProductId, true).Return(true, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch, protectedBranchList)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})

	t.Run("Should return error when GetProductByName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{}, errors.New("product not found"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errGetProduct, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should return error when GetEngagements fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, errors.New("api error"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errGetEngagements, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should return error when UpdateEngagementEndDate fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: pastDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().UpdateEngagementEndDate(testEngagementId, testProductId, false).Return(false, errors.New("update failed"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errUpdateEndDate, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should return error when CreateEngagement fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)
		mockService.EXPECT().CreateEngagement(testProjectName, testBranch, testProductId, false).Return(0, errors.New("create failed"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errCreateEngagement, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should skip non-matching engagements and create a new one", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		otherEngagement := defectdojo.Engagement{
			Id:        99,
			Name:      "ScopeGuardian-other-branch",
			TargetEnd: futureDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{otherEngagement}, nil)
		mockService.EXPECT().CreateEngagement(testProjectName, testBranch, testProductId, false).Return(testEngagementId, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})
}

func TestGetActiveFindings(t *testing.T) {
	t.Run("Should return active findings for an existing engagement", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: futureDate(),
		}

		ddFindings := []defectdojo.Finding{
			{Id: 1, Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42},
			{Id: 2, Title: "CVE-2021-1234", Severity: "Critical", FilePath: "/app/package.json", Line: 0},
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return(ddFindings, nil).AnyTimes()

		findings, err := GetActiveFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Len(t, findings, 2)
		assert.Equal(t, ddFindings, findings)
	})

	t.Run("Should return empty slice when no active findings exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: futureDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return([]defectdojo.Finding{}, nil).AnyTimes()

		findings, err := GetActiveFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Empty(t, findings)
	})

	t.Run("Should return error when GetProductByName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{}, errors.New("product not found"))

		findings, err := GetActiveFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})

	t.Run("Should return error when no matching engagement exists for branch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)

		findings, err := GetActiveFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errEngagementNotFound, err.Error())
		assert.Nil(t, findings)
	})

	t.Run("Should return error when GetFindings fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: futureDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return(nil, errors.New("api error"))

		findings, err := GetActiveFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errGetFindings, err.Error())
		assert.Nil(t, findings)
	})
}

func TestFilterByActiveFindings(t *testing.T) {
	// localFinding is a helper that builds a models.Finding with its Hash pre-computed,
	// matching what a scanner's LoadFindings would produce.
	localFinding := func(vulnId, severity, name, sinkFile string, sinkLine int, recommendation string) models.Finding {
		f := models.Finding{
			Name:           name,
			VulnId:         vulnId,
			Severity:       severity,
			SinkFile:       sinkFile,
			SinkLine:       sinkLine,
			Recommendation: recommendation,
		}
		f.Hash = models.ComputeFindingHash(vulnId, severity, sinkFile, sinkLine, recommendation)
		return f
	}

	t.Run("Should keep local findings whose hash matches a DD active finding", func(t *testing.T) {
		local := []models.Finding{
			localFinding("", "HIGH", "SQL Injection", "src/db.go", 42, "Use parameterized queries"),
			localFinding("", "MEDIUM", "XSS", "src/handler.go", 17, "Sanitize output"),
		}
		active := []defectdojo.Finding{
			{Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42, Mitigation: "Use parameterized queries"},
			{Title: "XSS", Severity: "Medium", FilePath: "src/handler.go", Line: 17, Mitigation: "Sanitize output"},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 2)
	})

	t.Run("Should remove local findings whose hash has no DD counterpart", func(t *testing.T) {
		local := []models.Finding{
			localFinding("", "HIGH", "SQL Injection", "src/db.go", 42, "Use parameterized queries"),
			localFinding("", "MEDIUM", "XSS", "src/handler.go", 17, "Sanitize output"),
		}
		active := []defectdojo.Finding{
			{Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42, Mitigation: "Use parameterized queries"},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
		assert.Equal(t, "SQL Injection", filtered[0].Name)
	})

	t.Run("Should match Grype findings via CVE in vulnerability_ids", func(t *testing.T) {
		local := []models.Finding{
			localFinding("CVE-2021-1234", "HIGH", "test-package 1.0.0", "/app/package.json", 0, "Upgrade to 2.0.0"),
			localFinding("CVE-2021-5678", "HIGH", "another-package 2.0.0", "/app/package.json", 0, "Upgrade to 3.0.0"),
		}
		// DD suppressed CVE-2021-5678 (false positive), only CVE-2021-1234 is active.
		// DefectDojo's Anchore Grype parser stores vulnerability.id in vulnerability_ids.
		active := []defectdojo.Finding{
			{
				Title:      "CVE-2021-1234",
				Severity:   "High",
				FilePath:   "/app/package.json",
				Line:       0,
				Mitigation: "Upgrade to 2.0.0",
				VulnerabilityIds: []defectdojo.VulnerabilityId{
					{VulnerabilityId: "CVE-2021-1234"},
				},
			},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
		assert.Equal(t, "test-package 1.0.0", filtered[0].Name)
		assert.Equal(t, "CVE-2021-1234", filtered[0].VulnId)
	})

	t.Run("Should match Grype findings using VulnerabilityIds array from DD", func(t *testing.T) {
		// DefectDojo also stores the CVE ID in the vulnerability_ids array; the hash must
		// be computable from that field alone so matching works regardless of how DD formats
		// the title vs the vulnerability_ids entries.
		local := []models.Finding{
			localFinding("CVE-2021-9999", "CRITICAL", "vuln-pkg 3.0.0", "/app/go.sum", 0, "Upgrade to 4.0.0"),
		}
		active := []defectdojo.Finding{
			{
				Title:    "CVE-2021-9999",
				Severity: "Critical",
				FilePath: "/app/go.sum",
				Line:     0,
				Mitigation: "Upgrade to 4.0.0",
				VulnerabilityIds: []defectdojo.VulnerabilityId{
					{VulnerabilityId: "CVE-2021-9999", Url: "https://nvd.nist.gov/vuln/detail/CVE-2021-9999"},
				},
			},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
		assert.Equal(t, "CVE-2021-9999", filtered[0].VulnId)
	})

	t.Run("Should match OpenGrep findings via check_id in vulnerability_ids after injection", func(t *testing.T) {
		// OpenGrep: enrichOpenGrepResults injects check_id into extra.metadata.cve before
		// upload; DefectDojo's Semgrep parser stores it in vulnerability_ids. Hash uses
		// check_id + severity + file + line with empty recommendation because DD stores
		// extra.message in description, not mitigation, for Semgrep-format findings.
		local := []models.Finding{
			localFinding("go.lang.security.injection.sql", "HIGH", "go.lang.security.injection.sql", "src/db.go", 10, ""),
		}
		active := []defectdojo.Finding{
			{
				Title:      "go.lang.security.injection.sql",
				Severity:   "High",
				FilePath:   "src/db.go",
				Line:       10,
				Mitigation: "",
				VulnerabilityIds: []defectdojo.VulnerabilityId{
					{VulnerabilityId: "go.lang.security.injection.sql"},
				},
			},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
	})

	t.Run("Should match KICS findings despite category-prefixed DD title", func(t *testing.T) {
		local := []models.Finding{
			localFinding("", "HIGH", "Missing User Instruction", "../../go/src/ScopeGuardian/Dockerfile", 81, "The 'Dockerfile' should contain the 'USER' instruction"),
		}
		// DD prefixes KICS finding titles with the rule category; hash matching ignores the title.
		active := []defectdojo.Finding{
			{
				Title:      "Build Process: Missing User Instruction",
				Severity:   "High",
				FilePath:   "../../go/src/ScopeGuardian/Dockerfile",
				Line:       81,
				Mitigation: "The 'Dockerfile' should contain the 'USER' instruction",
			},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
		assert.Equal(t, "Missing User Instruction", filtered[0].Name)
	})

	t.Run("Should not match when hash attributes differ", func(t *testing.T) {
		local := []models.Finding{
			localFinding("", "LOW", "Healthcheck Not Set", "Dockerfile", 10, "Add HEALTHCHECK"),
		}
		active := []defectdojo.Finding{
			{Title: "Build Process: Missing User Instruction", Severity: "High", FilePath: "Dockerfile", Line: 81, Mitigation: "Add USER instruction"},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Empty(t, filtered)
	})

	t.Run("Should be case-insensitive for severity when matching", func(t *testing.T) {
		local := []models.Finding{
			// Scanner emits severity in all-caps; DD returns title-case.
			localFinding("", "HIGH", "Vuln", "src/main.go", 5, "Fix it"),
		}
		active := []defectdojo.Finding{
			{Title: "Vuln", Severity: "High", FilePath: "src/main.go", Line: 5, Mitigation: "Fix it"},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
	})

	t.Run("Should return empty slice when active findings list is empty", func(t *testing.T) {
		local := []models.Finding{
			localFinding("", "HIGH", "SQL Injection", "src/db.go", 42, "Use parameterized queries"),
		}

		filtered := FilterByActiveFindings(local, []defectdojo.Finding{})

		assert.Empty(t, filtered)
	})

	t.Run("Should return empty slice when local findings list is empty", func(t *testing.T) {
		active := []defectdojo.Finding{
			{Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42, Mitigation: "Use parameterized queries"},
		}

		filtered := FilterByActiveFindings([]models.Finding{}, active)

		assert.Empty(t, filtered)
	})

	t.Run("Should not match when local finding has no hash set (unhashed finding)", func(t *testing.T) {
		// A finding that was never passed through LoadFindings (no Hash) should never match.
		local := []models.Finding{
			{Name: "Orphan Finding", Severity: "HIGH", SinkFile: "src/db.go", SinkLine: 1},
		}
		active := []defectdojo.Finding{
			{Title: "Orphan Finding", Severity: "High", FilePath: "src/db.go", Line: 1, Mitigation: ""},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Empty(t, filtered)
	})
}
