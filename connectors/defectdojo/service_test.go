package defectdojo

import (
	"net/http"
	"testing"

	"scope-guardian/connectors/defectdojo/client"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	URL          = "http://localhost:8080"
	TOKEN        = "f2f75000d4ubeae880ce72c3fe83ea342543af27"
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
