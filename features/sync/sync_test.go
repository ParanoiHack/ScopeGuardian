package sync

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"ScopeGuardian/connectors/defectdojo"

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

func TestGetDefectDojoFindings(t *testing.T) {
	t.Run("Should return findings mapped to internal model", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		existingEngagement := defectdojo.Engagement{
			Id:        testEngagementId,
			Name:      expectedEngagementName(),
			TargetEnd: futureDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{existingEngagement}, nil)
		mockService.EXPECT().GetFindings(testEngagementId, 0, 100, []defectdojo.Finding{}).Return([]defectdojo.Finding{
			{Id: 1, Title: "SQL Injection", Severity: "Critical"},
			{Id: 2, Title: "XSS", Severity: "High"},
		}, nil)

		findings, err := GetDefectDojoFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Len(t, findings, 2)
		assert.Equal(t, "SQL Injection", findings[0].Name)
		assert.Equal(t, "Critical", findings[0].Severity)
		assert.Equal(t, "XSS", findings[1].Name)
		assert.Equal(t, "High", findings[1].Severity)
	})

	t.Run("Should return empty slice when no findings exist", func(t *testing.T) {
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

		findings, err := GetDefectDojoFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.Nil(t, err)
		assert.Empty(t, findings)
	})

	t.Run("Should return error when GetEngagementId fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{}, errors.New("product not found"))

		findings, err := GetDefectDojoFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
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

		findings, err := GetDefectDojoFindings(mockService, testProjectName, testBranch, noProtectedBranches)

		assert.NotNil(t, err)
		assert.Equal(t, errGetFindings, err.Error())
		assert.Nil(t, findings)
	})
}
