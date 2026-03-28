package sync

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"scope-guardian/connectors/defectdojo"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	testProjectName = "my-project"
	testBranch      = "main"
	testProductId   = 5
	testEngagementId = 42
)

func futureDate() string {
	return time.Now().AddDate(1, 0, 0).Format(defectdojo.DateFormat)
}

func pastDate() string {
	return time.Now().AddDate(-1, 0, 0).Format(defectdojo.DateFormat)
}

func expectedEngagementName() string {
	return fmt.Sprintf("%s-%s", defectdojo.EngagementPrefix, testBranch)
}

func TestGetEngagementId(t *testing.T) {
	t.Run("Should create engagement when none exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)
		mockService.EXPECT().CreateEngagement(testBranch, testProductId).Return(testEngagementId, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

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

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

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
		mockService.EXPECT().UpdateEngagementEndDate(testEngagementId, testProductId).Return(true, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})

	t.Run("Should return error when GetProductByName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{}, errors.New("product not found"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

		assert.NotNil(t, err)
		assert.Equal(t, errGetProduct, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should return error when GetEngagements fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, errors.New("api error"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

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
		mockService.EXPECT().UpdateEngagementEndDate(testEngagementId, testProductId).Return(false, errors.New("update failed"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

		assert.NotNil(t, err)
		assert.Equal(t, errUpdateEndDate, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should return error when CreateEngagement fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{}, nil)
		mockService.EXPECT().CreateEngagement(testBranch, testProductId).Return(0, errors.New("create failed"))

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

		assert.NotNil(t, err)
		assert.Equal(t, errCreateEngagement, err.Error())
		assert.Equal(t, 0, id)
	})

	t.Run("Should skip non-matching engagements and create a new one", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := defectdojo.NewMockDefectDojoService(ctrl)

		otherEngagement := defectdojo.Engagement{
			Id:        99,
			Name:      "scope-guardian-other-branch",
			TargetEnd: futureDate(),
		}

		mockService.EXPECT().GetProductByName(testProjectName).Return(defectdojo.Product{Id: testProductId}, nil)
		mockService.EXPECT().GetEngagements(uint(testProductId), 0, 100, []defectdojo.Engagement{}).Return([]defectdojo.Engagement{otherEngagement}, nil)
		mockService.EXPECT().CreateEngagement(testBranch, testProductId).Return(testEngagementId, nil)

		id, err := GetEngagementId(mockService, testProjectName, testBranch)

		assert.Nil(t, err)
		assert.Equal(t, testEngagementId, id)
	})
}
