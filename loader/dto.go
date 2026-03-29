package loader

type (
	// Config holds the top-level application configuration loaded from the TOML file.
	Config struct {
		Title             string
		Kics              Kics
		Grype             Grype
		ProtectedBranches []string `toml:"protected_branches"`
	}

	// Kics contains the configuration for the KICS infrastructure-as-code scanner.
	Kics struct {
		Path     string
		Platform string
	}

	// Grype contains the configuration for the Grype vulnerability scanner.
	// When present, it also triggers the Syft SBOM generation step.
	Grype struct {
		Path string
	}
)
