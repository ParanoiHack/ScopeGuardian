package defectdojo

import (
	"net/http"
	"os"
	"testing"

	"scope-guardian/connectors/defectdojo/client"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	URL          = "http://172.19.0.7:8080"
	TOKEN        = "9218806a28dd10aa8a3cb6641e2e9f079d3f464e"
	PROJECT_NAME = "WebGoat"
)

func TestGetProductByName(t *testing.T) {
	gomockController := gomock.NewController(t)

	t.Run("Should retrieve product by name", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"count": 1,
				"results": [
					{
						"id": 1
					}
				]
			}
		`)

		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		product, err := service.GetProductByName(PROJECT_NAME)

		assert.Nil(t, err)
		assert.EqualValues(t, product.Id, 1)
	})

	t.Run("Should retrieve more than one product for a given name", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"count": 2,
				"results": [
					{
						"id": 1
					},
					{
						"id": 2
					}
				]
			}
		`)

		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		product, err := service.GetProductByName(PROJECT_NAME)

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Error(), errDuplicateProduct)
		assert.EqualValues(t, product.Id, 0)
	})

	t.Run("Should not retrieve product due to a lack of authorization", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"count": 1,
				"results": [
					{
						"id": 1
					}
				]
			}
		`)

		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 401).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		product, err := service.GetProductByName(PROJECT_NAME)

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Error(), errRetrieveProducts)
		assert.EqualValues(t, product.Id, 0)
	})

	t.Run("Should not retrieve product due to wrong json formating", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
				"count": 1,
				"results": [
					{
						"id": 1
					}
				]
			}
		`)

		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		product, err := service.GetProductByName(PROJECT_NAME)

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Error(), errUnmarshal)
		assert.EqualValues(t, product.Id, 0)
	})

	t.Run("Should not retrieve product due to wrong name", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"count": 0,
				"results": [
				]
			}
		`)

		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		product, err := service.GetProductByName(PROJECT_NAME)

		assert.NotNil(t, err)
		assert.EqualValues(t, err.Error(), errProductNotExist)
		assert.EqualValues(t, product.Id, 0)
	})
}

// func TestGetEngagements(t *testing.T) {
// 	// gomockController := gomock.NewController(t)

// 	t.Run("Should retrieve product engagements", func(t *testing.T) {
// 		// clientMock := client.NewMockClient(gomockController)

// 		// responseReturnMock := []byte(`
// 		// 	{
// 		// 		"count": 1,
// 		// 		"results": [
// 		// 			{
// 		// 				"id": 1
// 		// 			}
// 		// 		]
// 		// 	}
// 		// `)

// 		// clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
// 		// clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

// 		service := newDefectDojoService(client.NewClient(&http.Client{}), URL, TOKEN)

// 		engagements, err := service.GetEngagements(1, 0, 1, []Engagement{})

// 		assert.Nil(t, err)
// 		assert.EqualValues(t, len(engagements), 2)
// 	})
// }

// func TestCreateEngagement(t *testing.T) {
// 	// gomockController := gomock.NewController(t)

// 	t.Run("Should create an engagement", func(t *testing.T) {
// 		// clientMock := client.NewMockClient(gomockController)

// 		// responseReturnMock := []byte(`
// 		// 	{
// 		// 		"count": 1,
// 		// 		"results": [
// 		// 			{
// 		// 				"id": 1
// 		// 			}
// 		// 		]
// 		// 	}
// 		// `)

// 		// clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
// 		// clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

// 		service := newDefectDojoService(client.NewClient(&http.Client{}), URL, TOKEN)

// 		id, err := service.CreateEngagement("master", 4)

// 		fmt.Println(id)
// 		assert.Nil(t, err)
// 	})
// }

// func TestUpdateEngagementEndDate(t *testing.T) {
// 	// gomockController := gomock.NewController(t)

// 	t.Run("Should update engagement end date", func(t *testing.T) {
// 		// clientMock := client.NewMockClient(gomockController)

// 		// responseReturnMock := []byte(`
// 		// 	{
// 		// 		"count": 1,
// 		// 		"results": [
// 		// 			{
// 		// 				"id": 1
// 		// 			}
// 		// 		]
// 		// 	}
// 		// `)

// 		// clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
// 		// clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

// 		service := newDefectDojoService(client.NewClient(&http.Client{}), URL, TOKEN)

// 		ok, err := service.UpdateEngagementEndDate(13, 4)

// 		assert.Nil(t, err)
// 		assert.True(t, ok)
// 	})
// }

func TestImportScan(t *testing.T) {
	// gomockController := gomock.NewController(t)

	t.Run("Should import scan", func(t *testing.T) {
		// clientMock := client.NewMockClient(gomockController)

		// responseReturnMock := []byte(`
		// 	{
		// 		"count": 1,
		// 		"results": [
		// 			{
		// 				"id": 1
		// 			}
		// 		]
		// 	}
		// `)

		// clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
		// clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(client.NewClient(&http.Client{}), URL, TOKEN)

		var payload ScanPayload

		// payload.SeverityThreshold = "Info"
		// payload.Branch = "test"
		// payload.Tags = []string{"TEST"}
		// payload.GroupBy = "finding_title"
		// payload.FindingGroup = true
		// payload.FindingTag = true
		// payload.CloseOldFinding = true
		payload.EngagementId = 13
		// payload.Timestamp = time.Now().Format("2006-01-02")
		payload.ScanType = "KICS Scan"

		payload.File, _ = os.ReadFile("../../features/scans/kics/mocks/working_results/results/kics-results.json")

		ok, err := service.ImportScan(payload, "../../features/scans/kics/mocks/working_results/results/kics-results.json")

		assert.Nil(t, err)
		assert.True(t, ok)
	})
}
