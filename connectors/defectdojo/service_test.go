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
	URL          = "http://localhost:8080"
	TOKEN        = "9218806a28dd10ad8a3cb6641e2e9f079d3f464e"
	PROJECT_NAME = "Test"
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

func TestGetEngagements(t *testing.T) {
	gomockController := gomock.NewController(t)

	t.Run("Should retrieve product engagements in a recursive fashion", func(t *testing.T) {
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

		responseEndMock := []byte(`
			{
				"count": 1,
				"results": []		
			}
		`)

		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()
		gomock.InOrder(
			clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200),
			clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseEndMock, 200),
		)

		service := newDefectDojoService(clientMock, URL, TOKEN)

		engagements, err := service.GetEngagements(1, 0, 1, []Engagement{})

		assert.Nil(t, err)
		assert.EqualValues(t, 1, len(engagements))
	})

	t.Run("Should not retrieve engagement due to wrong HTTP status code", func(t *testing.T) {
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

		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()
		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 403)

		service := newDefectDojoService(clientMock, URL, TOKEN)

		engagements, err := service.GetEngagements(1, 0, 1, []Engagement{})

		assert.NotNil(t, err)
		assert.Equal(t, errRetrieveEngagements, err.Error())
		assert.EqualValues(t, 0, len(engagements))
	})

	t.Run("Should not retrieve engagement due to wrong JSON object", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"count": 1,
				"results": [
					{
						"id": 1
					}
				]
		`)

		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()
		clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200)

		service := newDefectDojoService(clientMock, URL, TOKEN)

		engagements, err := service.GetEngagements(1, 0, 1, []Engagement{})

		assert.NotNil(t, err)
		assert.Equal(t, errUnmarshal, err.Error())
		assert.EqualValues(t, 0, len(engagements))
	})

	t.Run("Should retrieve product engagements in a recursive fashion bis", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock_1 := []byte(`
			{
				"count": 2,
				"results": [
					{
						"id": 1
					}
				]
			}
		`)

		responseReturnMock_2 := []byte(`
			{
				"count": 2,
				"results": [
					{
						"id": 1
					}
				]
			}
		`)

		responseEndMock := []byte(`
			{
				"count": 2,
				"results": []		
			}
		`)

		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()
		gomock.InOrder(
			clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock_1, 200),
			clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock_2, 200),
			clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseEndMock, 200),
		)

		service := newDefectDojoService(clientMock, URL, TOKEN)

		engagements, err := service.GetEngagements(1, 0, 1, []Engagement{})

		assert.Nil(t, err)
		assert.EqualValues(t, 2, len(engagements))
	})

	t.Run("Should not retrieve product engagements cause not existing", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"count": 0,
				"results": [
				]
			}
		`)

		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()
		gomock.InOrder(
			clientMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(responseReturnMock, 200),
		)

		service := newDefectDojoService(clientMock, URL, TOKEN)

		engagements, err := service.GetEngagements(1, 0, 1, []Engagement{})

		assert.Nil(t, err)
		assert.EqualValues(t, 0, len(engagements))
	})
}

func TestCreateEngagement(t *testing.T) {
	gomockController := gomock.NewController(t)

	t.Run("Should create an engagement", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"id": 1
			}
		`)

		clientMock.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return(responseReturnMock, 201).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		id, err := service.CreateEngagement("master", 4)

		assert.Nil(t, err)
		assert.EqualValues(t, 1, id)
	})

	t.Run("Should not create an engagement due to Unmarshalling error", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"id": 1
		`)

		clientMock.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return(responseReturnMock, 201).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		id, err := service.CreateEngagement("master", 1)

		assert.NotNil(t, err)
		assert.EqualValues(t, errUnmarshal, err.Error())
		assert.EqualValues(t, 0, id)
	})

	t.Run("Should not create an engagement due to wrong HTTP status code", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`
			{
				"id": 1
		`)

		clientMock.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return(responseReturnMock, 403).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		id, err := service.CreateEngagement("master", 1)

		assert.NotNil(t, err)
		assert.EqualValues(t, errCreateEngagement, err.Error())
		assert.EqualValues(t, 0, id)
	})
}

func TestUpdateEngagementEndDate(t *testing.T) {
	gomockController := gomock.NewController(t)

	t.Run("Should update engagement end date", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`{}`)

		clientMock.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(responseReturnMock, 200).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		ok, err := service.UpdateEngagementEndDate(13, 4)

		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Should not update engagement due to wrong HTTP status code", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		responseReturnMock := []byte(`{}`)

		clientMock.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(responseReturnMock, 403).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		ok, err := service.UpdateEngagementEndDate(13, 4)

		assert.NotNil(t, err)
		assert.EqualValues(t, errUpdateEngagementEndDate, err.Error())
		assert.False(t, ok)
	})
}

func TestImportScan(t *testing.T) {
	gomockController := gomock.NewController(t)

	t.Run("Should import scan", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		clientMock.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(""), 201).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		ok, err := service.ImportScan(ScanPayload{}, "../../features/scans/kics/mocks/working_results/results/kics-results.json")

		assert.Nil(t, err)
		assert.True(t, ok)
	})

	t.Run("Should not import scan", func(t *testing.T) {
		clientMock := client.NewMockClient(gomockController)

		clientMock.EXPECT().Post(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte(""), 403).AnyTimes()
		clientMock.EXPECT().GetHeaders(gomock.Any()).Return(http.Header{}).AnyTimes()

		service := newDefectDojoService(clientMock, URL, TOKEN)

		ok, err := service.ImportScan(ScanPayload{}, "../../features/scans/kics/mocks/working_results/results/kics-results.json")

		assert.NotNil(t, err)
		assert.EqualValues(t, errImportScan, err.Error())
		assert.False(t, ok)
	})
}

func TestCreateMultipartFromScanPayload(t *testing.T) {
	t.Run("Should generate expected multipart body", func(t *testing.T) {
		var payload ScanPayload

		payload.SeverityThreshold = "severityThreshold"
		payload.Branch = "branch"
		payload.Tags = []string{"IACSTEngineTag"}
		payload.GroupBy = "groupByProperty"
		payload.FindingGroup = true
		payload.FindingTag = true
		payload.ScanType = "scanType"
		payload.EngagementId = 12
		payload.CloseOldFinding = true
		payload.File, _ = os.ReadFile("../../features/scans/kics/mocks/working_results/results/kics-results.json")

		_, _, err := createMultipartFromScanPayload(payload, "test.json")

		assert.Nil(t, err)
	})
}
