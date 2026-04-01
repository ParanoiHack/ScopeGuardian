package grype

type GrypeFix struct {
	Versions []string `json:"versions"`
	State    string   `json:"state"`
}

type GrypeVulnerability struct {
	ID          string    `json:"id"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Fix         GrypeFix  `json:"fix"`
}

type GrypeArtifactLocation struct {
	Path string `json:"path"`
}

type GrypeArtifact struct {
	Name      string                  `json:"name"`
	Version   string                  `json:"version"`
	Type      string                  `json:"type"`
	Locations []GrypeArtifactLocation `json:"locations"`
}

type GrypeMatch struct {
	Vulnerability GrypeVulnerability `json:"vulnerability"`
	Artifact      GrypeArtifact      `json:"artifact"`
}

type GrypeResults struct {
	Matches []GrypeMatch `json:"matches"`
}
