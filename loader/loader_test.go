package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Should load configuration file", func(t *testing.T) {
		config, err := Load("./mocks/config.toml")

		assert.Nil(t, err)
		assert.EqualValues(t, "Scope-guardian configuration file", config.Title)
		assert.EqualValues(t, "./", config.Path)
		assert.NotNil(t, config.Kics)
		assert.EqualValues(t, "terraform", config.Kics.Platform)
		assert.EqualValues(t, []string{"a227ec01-f97a-4084-91a4-47b350c1db54"}, config.Kics.ExcludeQueries)
	})

	t.Run("Should load configuration file with grype section", func(t *testing.T) {
		config, err := Load("./mocks/config_with_grype.toml")

		assert.Nil(t, err)
		assert.EqualValues(t, "Scope-guardian configuration file", config.Title)
		assert.EqualValues(t, "./", config.Path)
		assert.NotNil(t, config.Grype)
		assert.Nil(t, config.Kics)
		assert.EqualValues(t, "not-fixed,unknown,wont-fix", config.Grype.IgnoreStates)
		assert.EqualValues(t, []string{"/etc/**"}, config.Grype.Exclude)
	})

	t.Run("Should load configuration file without engine", func(t *testing.T) {
		config, err := Load("./mocks/config_no_engine.toml")

		assert.Nil(t, err)
		assert.EqualValues(t, "Scope-guardian configuration file", config.Title)
		assert.Nil(t, config.Kics)
		assert.Nil(t, config.Grype)
		assert.Nil(t, config.Opengrep)
	})

	t.Run("Should load configuration file with opengrep section", func(t *testing.T) {
		config, err := Load("./mocks/config_with_opengrep.toml")

		assert.Nil(t, err)
		assert.EqualValues(t, "Scope-guardian configuration file", config.Title)
		assert.EqualValues(t, "./", config.Path)
		assert.NotNil(t, config.Opengrep)
		assert.Nil(t, config.Kics)
		assert.Nil(t, config.Grype)
		assert.EqualValues(t, []string{"tests/**", "vendor/**"}, config.Opengrep.Exclude)
		assert.EqualValues(t, []string{"python.lang.correctness"}, config.Opengrep.ExcludeRule)
	})

	t.Run("Should load configuration file with proxy section", func(t *testing.T) {
		config, err := Load("./mocks/config_with_proxy.toml")

		assert.Nil(t, err)
		assert.NotNil(t, config.Proxy)
		assert.EqualValues(t, "http://proxy.example.com:3128", config.Proxy.HttpProxy)
		assert.EqualValues(t, "http://proxy.example.com:3128", config.Proxy.HttpsProxy)
		assert.EqualValues(t, "localhost,127.0.0.1", config.Proxy.NoProxy)
		assert.EqualValues(t, "/etc/ssl/certs/burp-ca.pem", config.Proxy.SslCertFile)
	})

	t.Run("Should not load configuration file cause wrong pathname", func(t *testing.T) {
		config, err := Load("./mocks/does_not_exist.toml")

		assert.NotNil(t, err)
		assert.EqualValues(t, Config{}, config)
		assert.EqualValues(t, errFileNotFound, err.Error())
	})

	t.Run("Should not load configuration file cause corruption", func(t *testing.T) {
		config, err := Load("./mocks/config_corrupted.toml")

		assert.NotNil(t, err)
		assert.EqualValues(t, Config{}, config)
		assert.EqualValues(t, errDecodingToml, err.Error())
	})
}
