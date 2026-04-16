package sync

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"ScopeGuardian/connectors/defectdojo"
	"ScopeGuardian/domains/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return(ddFindings, nil)

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
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return([]defectdojo.Finding{}, nil)

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
	t.Run("Should keep local findings that match DD active findings by Name and file", func(t *testing.T) {
		local := []models.Finding{
			{Name: "SQL Injection", SinkFile: "src/db.go", SinkLine: 42, Severity: "High"},
			{Name: "XSS", SinkFile: "src/handler.go", SinkLine: 17, Severity: "Medium"},
		}
		active := []defectdojo.Finding{
			{Title: "SQL Injection", FilePath: "src/db.go", Line: 42},
			{Title: "XSS", FilePath: "src/handler.go", Line: 17},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 2)
	})

	t.Run("Should remove local findings not present in DD active findings", func(t *testing.T) {
		local := []models.Finding{
			{Name: "SQL Injection", SinkFile: "src/db.go", SinkLine: 42, Severity: "High"},
			{Name: "XSS", SinkFile: "src/handler.go", SinkLine: 17, Severity: "Medium"},
		}
		active := []defectdojo.Finding{
			{Title: "SQL Injection", FilePath: "src/db.go", Line: 42},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
		assert.Equal(t, "SQL Injection", filtered[0].Name)
	})

	t.Run("Should use VulnId as title key for Grype findings", func(t *testing.T) {
		local := []models.Finding{
			{Name: "test-package 1.0.0", VulnId: "CVE-2021-1234", SinkFile: "/app/package.json", SinkLine: 0, Engine: "SCA"},
			{Name: "another-package 2.0.0", VulnId: "CVE-2021-5678", SinkFile: "/app/package.json", SinkLine: 0, Engine: "SCA"},
		}
		// DD suppressed CVE-2021-5678 (false positive), only CVE-2021-1234 is active
		active := []defectdojo.Finding{
			{Title: "CVE-2021-1234", FilePath: "/app/package.json", Line: 0},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
		assert.Equal(t, "test-package 1.0.0", filtered[0].Name)
		assert.Equal(t, "CVE-2021-1234", filtered[0].VulnId)
	})

	t.Run("Should fall back to Name when VulnId is empty", func(t *testing.T) {
		local := []models.Finding{
			{Name: "go.lang.security.injection.sql", SinkFile: "src/db.go", SinkLine: 10},
		}
		active := []defectdojo.Finding{
			{Title: "go.lang.security.injection.sql", FilePath: "src/db.go", Line: 10},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
	})

	t.Run("Should return empty slice when no local findings match active findings", func(t *testing.T) {
		local := []models.Finding{
			{Name: "SQL Injection", SinkFile: "src/db.go", SinkLine: 42},
		}
		active := []defectdojo.Finding{}

		filtered := FilterByActiveFindings(local, active)

		assert.Empty(t, filtered)
	})

	t.Run("Should return empty slice when local findings is empty", func(t *testing.T) {
		active := []defectdojo.Finding{
			{Title: "SQL Injection", FilePath: "src/db.go", Line: 42},
		}

		filtered := FilterByActiveFindings([]models.Finding{}, active)

		assert.Empty(t, filtered)
	})

	t.Run("Should match even when file path differs, using title-only matching", func(t *testing.T) {
		local := []models.Finding{
			{Name: "SQL Injection", SinkFile: "src/db.go", SinkLine: 42},
		}
		active := []defectdojo.Finding{
			{Title: "SQL Injection", FilePath: "src/other.go", Line: 42},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
	})

	t.Run("Should match even when line differs, using title-only matching", func(t *testing.T) {
		local := []models.Finding{
			{Name: "SQL Injection", SinkFile: "src/db.go", SinkLine: 42},
		}
		active := []defectdojo.Finding{
			{Title: "SQL Injection", FilePath: "src/db.go", Line: 99},
		}

		filtered := FilterByActiveFindings(local, active)

		assert.Len(t, filtered, 1)
	})
}
