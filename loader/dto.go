package loader

type (
	// Config holds the top-level application configuration loaded from the TOML file.
	Config struct {
		Title             string
		Path              string
		Kics              *Kics
		Grype             *Grype
		Opengrep          *Opengrep
		Proxy             *Proxy
		ProtectedBranches []string `toml:"protected_branches"`
	}

	// Kics contains the configuration for the KICS infrastructure-as-code scanner.
	Kics struct {
		Platform       string
		ExcludeQueries []string `toml:"exclude_queries"`
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

	// Opengrep contains the configuration for the OpenGrep SAST scanner.
	Opengrep struct {
		Exclude     []string `toml:"exclude"`
		ExcludeRule []string `toml:"exclude_rule"`
	}

	// Proxy holds optional HTTP proxy settings forwarded to scanner sub-processes
	// as HTTP_PROXY, HTTPS_PROXY, NO_PROXY, SSL_CERT_FILE, and REQUESTS_CA_BUNDLE
	// environment variables. All fields are optional and are omitted from the child
	// environment when empty.
	Proxy struct {
		HttpProxy   string `toml:"http_proxy"`
		HttpsProxy  string `toml:"https_proxy"`
		NoProxy     string `toml:"no_proxy"`
		SslCertFile string `toml:"ssl_cert_file"`
	}
)

// ToEnv converts the Proxy configuration into a list of "KEY=VALUE" environment
// variable entries suitable for passing as extraEnv to exec.Wrap / exec.WrapAllowExitCodes.
// Both the uppercase (HTTP_PROXY) and lowercase (http_proxy) variants are included
// for maximum compatibility across tools. SSL_CERT_FILE is also emitted as
// REQUESTS_CA_BUNDLE so that Python-based tools (e.g. OpenGrep) honour the same
// certificate. Returns nil when the receiver is nil or all fields are empty.
func (p *Proxy) ToEnv() []string {
	if p == nil {
		return nil
	}

	var env []string
	if p.HttpProxy != "" {
		env = append(env, "HTTP_PROXY="+p.HttpProxy, "http_proxy="+p.HttpProxy)
	}
	if p.HttpsProxy != "" {
		env = append(env, "HTTPS_PROXY="+p.HttpsProxy, "https_proxy="+p.HttpsProxy)
	}
	if p.NoProxy != "" {
		env = append(env, "NO_PROXY="+p.NoProxy, "no_proxy="+p.NoProxy)
	}
	if p.SslCertFile != "" {
		env = append(env, "SSL_CERT_FILE="+p.SslCertFile, "REQUESTS_CA_BUNDLE="+p.SslCertFile)
	}
	return env
}
