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

// DefectDojoService defines the operations available against the DefectDojo API.
type DefectDojoService interface {
	GetProductByName(productName string) (Product, error)
	CreateEngagement(projectName string, branch string, productId int, protected bool) (int, error)
	GetEngagements(productId uint, offset int, limit int, engagements []Engagement) ([]Engagement, error)
	UpdateEngagementEndDate(engagementId, productId int, protected bool) (bool, error)
	ImportScan(payload ScanPayload, filename string) (bool, error)
	GetFindings(engagementId int, offset int, limit int, findings []Finding) ([]Finding, error)
	SetAccessToken(token string)
	SetURL(url string)
}

// DefectDojoServiceImpl is the concrete implementation of DefectDojoService.
type DefectDojoServiceImpl struct {
	client      client.Client
	accessToken string
	url         string
}

// newDefectDojoService creates a DefectDojoServiceImpl pre-configured with the given
// HTTP client, API base URL, and access token.
func newDefectDojoService(client client.Client, url string, accessToken string) DefectDojoService {
	return &DefectDojoServiceImpl{
		client:      client,
		accessToken: accessToken,
		url:         url,
	}
}

// newDefectDojoServiceWithoutParam creates a DefectDojoServiceImpl with only an HTTP client,
// leaving the URL and access token to be set later via SetURL and SetAccessToken.
func newDefectDojoServiceWithoutParam(client client.Client) DefectDojoService {
	return &DefectDojoServiceImpl{
		client: client,
	}
}

// SetURL sets the base URL used for all DefectDojo API requests.
func (s *DefectDojoServiceImpl) SetURL(url string) {
	s.url = url
}

// SetAccessToken sets the API access token used to authenticate against DefectDojo.
func (s *DefectDojoServiceImpl) SetAccessToken(token string) {
	s.accessToken = token
}

// GetProductByName fetches the DefectDojo product whose name matches productName.
// It returns an error if the product does not exist, cannot be retrieved, or if
// multiple products share the same name.
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

// GetEngagements retrieves all engagements for the given product using cursor-based
// pagination (offset/limit). It accumulates results recursively until all pages
// have been fetched and returns the complete slice.
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

// CreateEngagement creates a new CI engagement in DefectDojo for the given project and branch.
// If protected is true the engagement end date is set to one year from today; otherwise it is
// set to one week from today. It returns the ID of the newly created engagement.
func (s *DefectDojoServiceImpl) CreateEngagement(projectName string, branch string, productId int, protected bool) (int, error) {
	var payload EngagementPayload

	t := time.Now()

	payload.EngagementType = EngagementType
	payload.Tags = append(payload.Tags, []string{ScopeGuardianTag, branch}...)
	payload.Name = fmt.Sprintf("%s-%s", projectName, branch)
	payload.Description = fmt.Sprintf(EngagementDescription, branch)
	payload.Status = EngagementStatus
	payload.Branch = branch
	payload.DeduplicationOnEngagement = false
	payload.Lead = EngagementDefaultLead
	payload.Product = productId
	payload.TargetStart = t.Format(DateFormat)
	if protected {
		payload.TargetEnd = t.AddDate(1, 0, 0).Format(DateFormat)
	} else {
		payload.TargetEnd = t.AddDate(0, 0, 7).Format(DateFormat)
	}

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

// UpdateEngagementEndDate updates the target end date of the given engagement.
// If protected is true the new end date is one year from today; otherwise it is one week from today.
// It returns true on success or an error if the update fails.
func (s *DefectDojoServiceImpl) UpdateEngagementEndDate(engagementId, productId int, protected bool) (bool, error) {
	var payload EngagementPayload

	t := time.Now()

	payload.TargetStart = t.Format(DateFormat)
	if protected {
		payload.TargetEnd = t.AddDate(1, 0, 0).Format(DateFormat)
	} else {
		payload.TargetEnd = t.AddDate(0, 0, 7).Format(DateFormat)
	}
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

// ImportScan uploads a scan result file to DefectDojo using the given ScanPayload.
// filename is used as the multipart form file name. Returns true on success.
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

	if code < http.StatusOK || code >= http.StatusMultipleChoices {
		logger.Error(fmt.Sprintf(logErrorImportScan, code))
		return false, errors.New(errImportScan)
	}

	return true, nil
}

// createMultipartFromScanPayload serialises a ScanPayload into a multipart/form-data
// body. Struct fields are mapped to form keys via the "form" struct tag. The raw
// scan file is appended under the "file" key using filename. It returns the body
// bytes, the multipart boundary content-type string, and any error encountered.
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

// GetFindings retrieves all findings for the given engagement and product using
// cursor-based pagination (offset/limit). It accumulates results recursively until
// all pages have been fetched and returns the complete slice.
func (s *DefectDojoServiceImpl) GetFindings(engagementId int, offset int, limit int, findings []Finding) ([]Finding, error) {
	var res GetFindingsResponse

	body, code := s.client.Get(fmt.Sprintf(
		"%s%s%s", s.url, APIPrefix, fmt.Sprintf(GetFindingsPath, engagementId, offset, limit)), s.client.GetHeaders(s.accessToken))
	if code != http.StatusOK {
		logger.Error(fmt.Sprintf(logErrorRetrieveFindings, engagementId))
		return []Finding{}, errors.New(errRetrieveFindings)
	}

	err := json.Unmarshal(body, &res)
	if err != nil {
		logger.Error(fmt.Sprintf(logErrorDecodingToken, err.Error()))
		return []Finding{}, errors.New(errUnmarshal)
	}

	findings = append(findings, res.Results...)
	if res.Count-(offset+limit) >= 0 {
		return s.GetFindings(engagementId, offset+limit, limit, findings)
	}

	return findings, nil
}
