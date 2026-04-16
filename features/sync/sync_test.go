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

func TestMarkFindingsByDDFindings(t *testing.T) {
	// localFinding is a helper that builds a models.Finding with its Hash pre-computed,
	// matching what a scanner's LoadFindings would produce.
	localFinding := func(severity, name, sinkFile string, sinkLine int, recommendation string) models.Finding {
		f := models.Finding{
			Name:           name,
			Severity:       severity,
			SinkFile:       sinkFile,
			SinkLine:       sinkLine,
			Recommendation: recommendation,
		}
		f.Hash = models.ComputeFindingHash(severity, sinkFile, sinkLine, recommendation)
		return f
	}

	t.Run("Should mark finding as ACTIVE when DD finding is active and not duplicate", func(t *testing.T) {
		local := []models.Finding{
			localFinding("HIGH", "SQL Injection", "src/db.go", 42, "Use parameterized queries"),
		}
		ddFindings := []defectdojo.Finding{
			{Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42,
				Mitigation: "Use parameterized queries", Active: true, Duplicate: false},
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusActive, result[0].Status)
	})

	t.Run("Should mark finding as INACTIVE when DD finding is not active and not duplicate", func(t *testing.T) {
		local := []models.Finding{
			localFinding("MEDIUM", "XSS", "src/handler.go", 17, "Sanitize output"),
		}
		ddFindings := []defectdojo.Finding{
			{Title: "XSS", Severity: "Medium", FilePath: "src/handler.go", Line: 17,
				Mitigation: "Sanitize output", Active: false, Duplicate: false},
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusInactive, result[0].Status)
	})

	t.Run("Should mark finding as DUPLICATE when DD finding has duplicate=true", func(t *testing.T) {
		local := []models.Finding{
			localFinding("CRITICAL", "RCE", "src/exec.go", 5, "Avoid exec"),
		}
		ddFindings := []defectdojo.Finding{
			{Title: "RCE", Severity: "Critical", FilePath: "src/exec.go", Line: 5,
				Mitigation: "Avoid exec", Active: false, Duplicate: true},
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusDuplicate, result[0].Status)
	})

	t.Run("Should default to ACTIVE when local finding has no DD counterpart", func(t *testing.T) {
		local := []models.Finding{
			localFinding("HIGH", "New Finding", "src/new.go", 1, "Fix it"),
		}

		result := MarkFindingsByDDFindings(local, []defectdojo.Finding{})

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusActive, result[0].Status)
	})

	t.Run("Should return all findings (no filtering)", func(t *testing.T) {
		local := []models.Finding{
			localFinding("HIGH", "SQL Injection", "src/db.go", 42, "Use parameterized queries"),
			localFinding("MEDIUM", "XSS", "src/handler.go", 17, "Sanitize output"),
			localFinding("LOW", "Orphan", "src/orphan.go", 99, "Fix it"),
		}
		ddFindings := []defectdojo.Finding{
			{Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42,
				Mitigation: "Use parameterized queries", Active: true, Duplicate: false},
			{Title: "XSS", Severity: "Medium", FilePath: "src/handler.go", Line: 17,
				Mitigation: "Sanitize output", Active: false, Duplicate: true},
			// Orphan finding is not in DD — defaults to ACTIVE.
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 3)
		assert.Equal(t, models.FindingStatusActive, result[0].Status)
		assert.Equal(t, models.FindingStatusDuplicate, result[1].Status)
		assert.Equal(t, models.FindingStatusActive, result[2].Status)
	})

	t.Run("Should prioritise duplicate=true over active=false for INACTIVE/DUPLICATE", func(t *testing.T) {
		local := []models.Finding{
			localFinding("HIGH", "Vuln", "src/main.go", 5, "Fix it"),
		}
		ddFindings := []defectdojo.Finding{
			{Title: "Vuln", Severity: "High", FilePath: "src/main.go", Line: 5,
				Mitigation: "Fix it", Active: false, Duplicate: true},
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusDuplicate, result[0].Status)
	})

	t.Run("Should match OpenGrep findings via UniqueIdFromTool and mark correctly", func(t *testing.T) {
		local := []models.Finding{
			localFinding("HIGH", "go.lang.security.injection.sql", "src/db.go", 10, ""),
		}
		ddFindings := []defectdojo.Finding{
			{
				Title:            "go.lang.security.injection.sql",
				Severity:         "High",
				FilePath:         "src/db.go",
				Line:             10,
				Mitigation:       "",
				UniqueIdFromTool: local[0].Hash,
				Active:           true,
				Duplicate:        false,
			},
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusActive, result[0].Status)
	})

	t.Run("Should mark KICS finding correctly despite category-prefixed DD title", func(t *testing.T) {
		local := []models.Finding{
			localFinding("HIGH", "Missing User Instruction", "../../go/src/ScopeGuardian/Dockerfile", 81,
				"The 'Dockerfile' should contain the 'USER' instruction"),
		}
		ddFindings := []defectdojo.Finding{
			{
				Title:      "Build Process: Missing User Instruction",
				Severity:   "High",
				FilePath:   "../../go/src/ScopeGuardian/Dockerfile",
				Line:       81,
				Mitigation: "The 'Dockerfile' should contain the 'USER' instruction",
				Active:     true,
				Duplicate:  false,
			},
		}

		result := MarkFindingsByDDFindings(local, ddFindings)

		assert.Len(t, result, 1)
		assert.Equal(t, models.FindingStatusActive, result[0].Status)
		assert.Equal(t, "Missing User Instruction", result[0].Name)
	})

	t.Run("Should return empty slice when local findings list is empty", func(t *testing.T) {
		ddFindings := []defectdojo.Finding{
			{Title: "SQL Injection", Severity: "High", FilePath: "src/db.go", Line: 42,
				Mitigation: "Fix it", Active: true, Duplicate: false},
		}

		result := MarkFindingsByDDFindings([]models.Finding{}, ddFindings)

		assert.Empty(t, result)
	})
}

func TestGetEngagementFindings(t *testing.T) {
	t.Run("Should poll then return all findings (active + inactive + duplicate)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: futureDate(),
		}

		// pollFindings uses GetFindings; return the same slice twice to trigger
		// the "stable count" exit condition quickly.
		pollFindings := []defectdojo.Finding{
			{Id: 1, Severity: "High", Active: true, Duplicate: false},
		}
		allFindings := []defectdojo.Finding{
			{Id: 1, Severity: "High", Active: true, Duplicate: false},
			{Id: 2, Severity: "Medium", Active: false, Duplicate: true},
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return(pollFindings, nil).AnyTimes()
		mockService.EXPECT().GetAllEngagementFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return(allFindings, nil)

		findings, err := GetEngagementFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Len(t, findings, 2)
	})

	t.Run("Should return error when no matching engagement exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)

		findings, err := GetEngagementFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errEngagementNotFound, err.Error())
		assert.Nil(t, findings)
	})

	t.Run("Should return error when GetProductByName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{}, errors.New("not found"))

		findings, err := GetEngagementFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})

	t.Run("Should return error when polling GetFindings fails", func(t *testing.T) {
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

		findings, err := GetEngagementFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errGetFindings, err.Error())
		assert.Nil(t, findings)
	})

	t.Run("Should return error when GetAllEngagementFindings fails", func(t *testing.T) {
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
		mockService.EXPECT().GetAllEngagementFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return(nil, errors.New("api error"))

		findings, err := GetEngagementFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Nil(t, findings)
	})
}
