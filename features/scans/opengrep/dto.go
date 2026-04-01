package opengrep

type OpenGrepResultLocation struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

type OpenGrepExtra struct {
	Message  string            `json:"message"`
	Severity string            `json:"severity"`
	Metadata map[string]interface{} `json:"metadata"`
}

type OpenGrepResult struct {
	CheckId string               `json:"check_id"`
	Path    string               `json:"path"`
	Start   OpenGrepResultLocation `json:"start"`
	Extra   OpenGrepExtra        `json:"extra"`
}

type OpenGrepResults struct {
	Results []OpenGrepResult `json:"results"`
}
