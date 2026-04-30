package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestConfigureTrustedProxiesNilIgnoresForwardedFor(t *testing.T) {
	t.Setenv("GIN_TRUSTED_PROXIES", "")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if !assert.NoError(t, configureTrustedProxies(r)) {
		return
	}
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.10:5555"
	req.Header.Set("X-Forwarded-For", "203.0.113.7")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "192.0.2.10", rec.Body.String())
}

func TestConfigureTrustedProxiesInvalidEnvReturnsError(t *testing.T) {
	t.Setenv("GIN_TRUSTED_PROXIES", "not-a-valid-cidr!!!")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	err := configureTrustedProxies(r)
	assert.Error(t, err)
}

func TestConfigureTrustedProxiesAllowsForwardedFromTrustedPeer(t *testing.T) {
	t.Setenv("GIN_TRUSTED_PROXIES", "192.0.2.10/32")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if !assert.NoError(t, configureTrustedProxies(r)) {
		return
	}
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.10:5555"
	req.Header.Set("X-Forwarded-For", "198.51.100.22")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "198.51.100.22", rec.Body.String())
}
