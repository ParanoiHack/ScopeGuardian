package defectdojo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"scope-guardian/connectors/defectdojo/client"
	"scope-guardian/logger"
	"strconv"
	"time"
)

type DefectDojoService interface {
	GetProductByName(productName string) (Product, error)
	CreateEngagement(branch string, productId int) (int, error)
	GetEngagements(productId uint, offset int, limit int, engagements []Engagement) ([]Engagement, error)
	UpdateEngagementEndDate(engagementId, productId int) (bool, error)
	ImportScan(payload ScanPayload, filename string) (bool, error)
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
		"%s%s%s%s", s.url, APIPrefix, GetProductsPath, productName), s.client.GetHeaders(s.accessToken))

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

func (s *DefectDojoServiceImpl) GetEngagements(productId uint, offset int, limit int, engagements []Engagement) ([]Engagement, error) {
	var res GetEngagementsResponse

	body, code := s.client.Get(fmt.Sprintf(
		"%s%s%s", s.url, APIPrefix, fmt.Sprintf(GetEngagementsPath, productId, offset, limit)), s.client.GetHeaders(s.accessToken))
	if code != http.StatusOK {
		logger.Error(fmt.Sprintf(logErrorRetrieveEngagements, productId))
		return []Engagement{}, errors.New(errRetrieveEngagements)
	}

	err := json.Unmarshal(body, &res)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorDecodingToken, err.Error()))
		return []Engagement{}, errors.New(errUnmarshal)
	}

	engagements = append(engagements, res.Results...)
	if res.Count-(offset+limit) >= 0 {
		return s.GetEngagements(productId, offset+limit, limit, engagements)
	}

	return engagements, nil
}

func (s *DefectDojoServiceImpl) CreateEngagement(branch string, productId int) (int, error) {
	var payload EngagementPayload

	t := time.Now()

	payload.EngagementType = EngagementType
	payload.Tags = append(payload.Tags, []string{ScopeGuardianTag, branch}...)
	payload.Name = fmt.Sprintf("%s-%s", EngagementPrefix, branch)
	payload.Description = fmt.Sprintf(EngagementDescription, branch)
	payload.Status = EngagementStatus
	payload.Branch = branch
	payload.DeduplicationOnEngagement = false
	payload.Lead = EngagementDefaultLead
	payload.Product = productId
	payload.TargetStart = t.Format(DateFormat)
	payload.TargetEnd = t.AddDate(1, 0, 0).Format(DateFormat)

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error(logErrorEncodingStruct)
		return 0, err
	}

	var res CreateEngagementResponse

	body, code := s.client.Post(fmt.Sprintf(
		"%s%s%s", s.url, APIPrefix, CreateEngagementPath), data, s.client.GetHeaders(s.accessToken))
	if code != http.StatusCreated {
		logger.Error(logErrorCreateEngagement)
		return 0, errors.New(errCreateEngagement)
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorDecodingToken, err.Error()))
		return 0, errors.New(errUnmarshal)
	}

	return res.Id, nil
}

func (s *DefectDojoServiceImpl) UpdateEngagementEndDate(engagementId, productId int) (bool, error) {
	var payload EngagementPayload

	t := time.Now()

	payload.TargetStart = t.Format(DateFormat)
	payload.TargetEnd = t.AddDate(1, 0, 0).Format(DateFormat)
	payload.Status = EngagementStatus
	payload.Product = productId

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error(logErrorEncodingStruct)
		return false, err
	}

	_, code := s.client.Put(fmt.Sprintf(
		"%s%s%s", s.url, APIPrefix, fmt.Sprintf(UpdateEngagementPath, engagementId)), data, s.client.GetHeaders(s.accessToken))
	if code != http.StatusOK {
		logger.Error(fmt.Sprintf(logErrorUpdateEngagementEndDate, engagementId))
		return false, errors.New(errUpdateEngagementEndDate)
	}

	return true, nil
}

func (s *DefectDojoServiceImpl) ImportScan(payload ScanPayload, filename string) (bool, error) {
	body, boundary, err := createMultipartFromScanPayload(payload, filename)
	if err != nil {
		logger.Error(logErrorCreateMultipartRequest)
		return false, err
	}

	headers := s.client.GetHeaders(s.accessToken)
	headers.Set(client.ContentTypeKey, boundary)

	_, code := s.client.Post(fmt.Sprintf(
		"%s%s%s", s.url, APIPrefix, ImportScanPath), body, headers)

	if code != http.StatusCreated {
		logger.Error(logErrorImportScan)
		return false, errors.New(errImportScan)
	}

	return true, nil
}

func createMultipartFromScanPayload(payload ScanPayload, filename string) ([]byte, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	payloadType := reflect.TypeOf(payload)
	payloadValue := reflect.ValueOf(payload)

	for i := 0; i < payloadType.NumField(); i++ {
		field := payloadType.Field(i)
		value := payloadValue.Field(i)

		formKey := field.Tag.Get(FormValuekey)
		if formKey == "" {
			continue
		}

		switch value.Kind() {
		case reflect.String:
			if err := writer.WriteField(formKey, value.String()); err != nil {
				logger.Error(fmt.Sprintf(logErrorReflection, formKey))
			}
		case reflect.Bool:
			valStr := strconv.FormatBool(value.Bool())
			if err := writer.WriteField(formKey, valStr); err != nil {
				logger.Error(fmt.Sprintf(logErrorReflection, formKey))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valStr := strconv.FormatInt(value.Int(), 10)
			if err := writer.WriteField(formKey, valStr); err != nil {
				logger.Error(fmt.Sprintf(logErrorReflection, formKey))
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				for _, item := range value.Interface().([]string) {
					if err := writer.WriteField(formKey, item); err != nil {
						logger.Error(fmt.Sprintf(logErrorReflection, formKey))
					}
				}
			}
		default:
			logger.Error(fmt.Sprintf(logErrorUnknownType, value.Kind()))
		}
	}

	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		logger.Error(logErrorAddFile)
	}

	_, err = filePart.Write(payload.File)
	if err != nil {
		return nil, "", errors.New(errWritingFile)
	}

	writer.Close()

	return body.Bytes(), writer.FormDataContentType(), nil
}

// get findings
