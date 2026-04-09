package parser

type (
	Args struct {
		Config      string
		ProjectName string
		Branch      string
		Sync        bool
		Quiet       bool
		Output      string
		Thresholds  []Threshold
	}

	Threshold struct {
		Severity string
		Value    int
	}
)
