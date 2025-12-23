package loader

type (
	Config struct {
		Title string
		Kics  Kics
	}

	Kics struct {
		Path     string
		Platform string
	}
)
