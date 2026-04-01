package opengrep

import "encoding/json"

type OpenGrepResultLocation struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

// StringOrSlice unmarshals a JSON value that may be either a plain string or
// an array of strings into a Go []string.
type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	// Try array first
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*s = slice
		return nil
	}
	// Fall back to plain string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = StringOrSlice{str}
	return nil
}

type OpenGrepMetadata struct {
	Cwe        StringOrSlice `json:"cwe"`
	Owasp      StringOrSlice `json:"owasp"`
	Impact     string        `json:"impact"`
	Confidence string        `json:"confidence"`
}

type OpenGrepExtra struct {
	Message  string           `json:"message"`
	Metadata OpenGrepMetadata `json:"metadata"`
}

type OpenGrepResult struct {
	CheckId string                 `json:"check_id"`
	Path    string                 `json:"path"`
	Start   OpenGrepResultLocation `json:"start"`
	Extra   OpenGrepExtra          `json:"extra"`
}

type OpenGrepResults struct {
	Results []OpenGrepResult `json:"results"`
}
