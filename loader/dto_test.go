package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyToEnv(t *testing.T) {
	t.Run("Should return nil for nil proxy", func(t *testing.T) {
		var p *Proxy
		assert.Nil(t, p.ToEnv())
	})

	t.Run("Should return nil for empty proxy", func(t *testing.T) {
		p := &Proxy{}
		assert.Nil(t, p.ToEnv())
	})

	t.Run("Should return HTTP_PROXY entries when http_proxy is set", func(t *testing.T) {
		p := &Proxy{HttpProxy: "http://proxy.example.com:3128"}
		env := p.ToEnv()
		assert.Contains(t, env, "HTTP_PROXY=http://proxy.example.com:3128")
		assert.Contains(t, env, "http_proxy=http://proxy.example.com:3128")
	})

	t.Run("Should return HTTPS_PROXY entries when https_proxy is set", func(t *testing.T) {
		p := &Proxy{HttpsProxy: "http://proxy.example.com:3128"}
		env := p.ToEnv()
		assert.Contains(t, env, "HTTPS_PROXY=http://proxy.example.com:3128")
		assert.Contains(t, env, "https_proxy=http://proxy.example.com:3128")
	})

	t.Run("Should return NO_PROXY entries when no_proxy is set", func(t *testing.T) {
		p := &Proxy{NoProxy: "localhost,127.0.0.1"}
		env := p.ToEnv()
		assert.Contains(t, env, "NO_PROXY=localhost,127.0.0.1")
		assert.Contains(t, env, "no_proxy=localhost,127.0.0.1")
	})

	t.Run("Should return all six entries when all fields are set", func(t *testing.T) {
		p := &Proxy{
			HttpProxy:  "http://proxy.example.com:3128",
			HttpsProxy: "https://proxy.example.com:3128",
			NoProxy:    "localhost",
		}
		env := p.ToEnv()
		assert.Len(t, env, 6)
	})
}
