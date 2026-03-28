package loader

type (
	// Config holds the top-level application configuration loaded from the TOML file.
	Config struct {
		Title string
		Kics  Kics
	}

	// Kics contains the configuration for the KICS infrastructure-as-code scanner.
	Kics struct {
		Path     string
		Platform string
	}
)
