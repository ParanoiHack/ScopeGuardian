package defectdojo

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMockDefectDojoService_GetProductByName_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	expected := Product{Id: 42}
	mock.EXPECT().GetProductByName("my-project").Return(expected, nil)

	product, err := mock.GetProductByName("my-project")
	assert.Nil(t, err)
	assert.Equal(t, 42, product.Id)
}

func TestMockDefectDojoService_GetProductByName_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().GetProductByName("unknown").Return(Product{}, errors.New(errProductNotExist))

	product, err := mock.GetProductByName("unknown")
	assert.NotNil(t, err)
	assert.Equal(t, 0, product.Id)
}

func TestMockDefectDojoService_GetEngagements_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	expected := []Engagement{{Id: 1, Name: "scope-guardian-main"}}
	mock.EXPECT().GetEngagements(uint(1), 0, 100, gomock.Any()).Return(expected, nil)

	engagements, err := mock.GetEngagements(1, 0, 100, []Engagement{})
	assert.Nil(t, err)
	assert.Len(t, engagements, 1)
	assert.Equal(t, 1, engagements[0].Id)
}

func TestMockDefectDojoService_GetEngagements_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().GetEngagements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]Engagement{}, errors.New(errRetrieveEngagements))

	engagements, err := mock.GetEngagements(1, 0, 100, []Engagement{})
	assert.NotNil(t, err)
	assert.Empty(t, engagements)
}

func TestMockDefectDojoService_CreateEngagement_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().CreateEngagement("my-project", "main", 1).Return(10, nil)

	id, err := mock.CreateEngagement("my-project", "main", 1)
	assert.Nil(t, err)
	assert.Equal(t, 10, id)
}

func TestMockDefectDojoService_CreateEngagement_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().CreateEngagement(gomock.Any(), gomock.Any(), gomock.Any()).Return(0, errors.New(errCreateEngagement))

	id, err := mock.CreateEngagement("my-project", "feature-branch", 1)
	assert.NotNil(t, err)
	assert.Equal(t, 0, id)
}

func TestMockDefectDojoService_UpdateEngagementEndDate_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().UpdateEngagementEndDate(5, 1).Return(true, nil)

	ok, err := mock.UpdateEngagementEndDate(5, 1)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestMockDefectDojoService_UpdateEngagementEndDate_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().UpdateEngagementEndDate(gomock.Any(), gomock.Any()).Return(false, errors.New(errUpdateEngagementEndDate))

	ok, err := mock.UpdateEngagementEndDate(5, 1)
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestMockDefectDojoService_ImportScan_Succeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().ImportScan(gomock.Any(), "results.json").Return(true, nil)

	ok, err := mock.ImportScan(ScanPayload{}, "results.json")
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestMockDefectDojoService_ImportScan_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().ImportScan(gomock.Any(), gomock.Any()).Return(false, errors.New(errImportScan))

	ok, err := mock.ImportScan(ScanPayload{}, "results.json")
	assert.NotNil(t, err)
	assert.False(t, ok)
}

func TestMockDefectDojoService_SetAccessToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().SetAccessToken("new-token")
	mock.SetAccessToken("new-token")
}

func TestMockDefectDojoService_SetURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockDefectDojoService(ctrl)

	mock.EXPECT().SetURL("http://new-url:8080")
	mock.SetURL("http://new-url:8080")
}
