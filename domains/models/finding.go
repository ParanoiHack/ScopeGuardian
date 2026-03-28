package models

type Finding struct {
	Engine         string
	Severity       string
	Name           string
	Cwe            string
	Description    string
	SinkFile       string
	SinkLine       int
	Recommendation string
}
