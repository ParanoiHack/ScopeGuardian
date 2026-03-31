package loader

type (
	// Config holds the top-level application configuration loaded from the TOML file.
	Config struct {
		Title             string
		Path              string
		Kics              *Kics
		Grype             *Grype
		ProtectedBranches []string `toml:"protected_branches"`
	}

	// Kics contains the configuration for the KICS infrastructure-as-code scanner.
	Kics struct {
		Platform string
	}

	// Grype contains the configuration for the Grype vulnerability scanner.
	// When present, it also triggers the Syft SBOM generation step.
	Grype struct {
		Exclude      []string `toml:"exclude"`
		IgnoreStates string   `toml:"ignore_states"`
		// TransitiveLibraries controls whether Syft resolves transitive Java
		// dependencies from Maven Central during SBOM generation. Disabled by
		// default because network resolution significantly increases scan time.
		TransitiveLibraries bool `toml:"transitive_libraries"`
	}
)
