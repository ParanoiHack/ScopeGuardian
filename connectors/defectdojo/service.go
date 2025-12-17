package defectdojo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"scope-guardian/connectors/defectdojo/client"
	"scope-guardian/logger"
)

type DefectDojoService interface {
	GetProductByName(productName string) (Product, error)
	SetAccessToken(token string)
	SetURL(url string)
}

type DefectDojoServiceImpl struct {
	client      client.Client
	accessToken string
	url         string
}

func newDefectDojoService(client client.Client, url string, accessToken string) DefectDojoService {
	return &DefectDojoServiceImpl{
		client:      client,
		accessToken: accessToken,
		url:         url,
	}
}

func newDefectDojoServiceWithoutParam(client client.Client) DefectDojoService {
	return &DefectDojoServiceImpl{
		client: client,
	}
}

func (s *DefectDojoServiceImpl) SetURL(url string) {
	s.url = url
}

func (s *DefectDojoServiceImpl) SetAccessToken(token string) {
	s.accessToken = token
}

func (s *DefectDojoServiceImpl) GetProductByName(productName string) (Product, error) {
	var res GetProductByNameResponse

	body, code := s.client.Get(fmt.Sprintf(
		"%s/%s/%s%s", s.url, APIPrefix, GetProductsPath, productName), s.client.GetHeaders(s.accessToken))

	if code != http.StatusOK {
		logger.Error(logErrorRetrieveProducts)
		return Product{}, errors.New(errRetrieveProducts)
	}

	err := json.Unmarshal(body, &res)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorDecodingToken, err.Error()))
		return Product{}, errors.New(errUnmarshal)
	}

	if res.Count == 1 {
		return res.Results[0], nil
	} else if res.Count > 1 {
		return Product{}, errors.New(errDuplicateProduct)
	}

	return Product{}, errors.New(errProductNotExist)
}
