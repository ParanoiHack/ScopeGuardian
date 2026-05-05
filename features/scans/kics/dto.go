package kics

type KicsFinding struct {
	QueryName   string     `json:"query_name"`
	QueryId     string     `json:"query_id"`
	QueryUrl    string     `json:"query_url"`
	Severity    string     `json:"severity"`
	Platform    string     `json:"platforme"`
	Cwe         string     `json:"cwe"`
	RiskScore   string     `json:"risk_score"`
	Description string     `json:"description"`
	Files       []KicsFile `json:"files"`
}

type KicsFile struct {
	FileName       string `json:"file_name"`
	SimilarityId   string `json:"similarity_id"`
	Line           int    `json:"line"`
	IssueType      string `json:"issue_type"`
	Recommendation string `json:"expected_value"`
}

type KicsResults struct {
	FilesScanned                       int           `json:"files_scanned"`
	LinesScanned                       int           `json:"lines_scanned"`
	FilesParsed                        int           `json:"files_parsed"`
	LinesIgnored                       int           `json:"lines_ingored"`
	LinesParsed                        int           `json:"lines_parsed"`
	FilesFailedToScan                  int           `json:"files_failed_to_scan"`
	QueriesTotal                       int           `json:"queries_total"`
	QueriesFailedToExecute             int           `json:"queries_failed_to_execute"`
	QueriesFailedToComputeSimilarityId int           `json:"queries_failed_to_conpute_similarity_id"`
	Queries                            []KicsFinding `json:"queries"`
}
