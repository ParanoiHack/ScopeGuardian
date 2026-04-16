package models

// Finding represents a single security finding produced by a scanner.
type Finding struct {
	Engine         string
	Severity       string
	Name           string
	VulnId         string
	Cwe            string
	Description    string
	SinkFile       string
	SinkLine       int
	Recommendation string
}
