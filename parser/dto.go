package parser

type (
	Args struct {
		Config    string
		Sync      bool
		Threshold *Threshold
	}

	Threshold struct {
		Severity string
		Value    int
	}
)
