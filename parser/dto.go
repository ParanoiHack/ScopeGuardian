package parser

type (
	Args struct {
		Config        string
		ProjectName   string
		Branch        string
		Sync          bool
		Quiet         bool
		Output        string
		Format        string
		Thresholds    []Threshold
		StatusFilters []string
	}

	Threshold struct {
		Severity string
		Value    int
	}
)
