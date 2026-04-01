package opengrep

type OpenGrepResultLocation struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

type OpenGrepMetadata struct {
	Cwe        []string `json:"cwe"`
	Owasp      []string `json:"owasp"`
	Impact     string   `json:"impact"`
	Confidence string   `json:"confidence"`
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
