package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefaultMaxRequestBodyBytes is used when MaxRequestBody is called with a
// non-positive limit or when the caller passes zero explicitly.
const DefaultMaxRequestBodyBytes int64 = 1 << 20 // 1 MiB

// MaxRequestBody returns middleware that caps how many bytes handlers may read
// from the request body for POST, PUT, and PATCH. Larger bodies yield HTTP 413
// (via http.MaxBytesReader) without buffering the entire payload in memory.
func MaxRequestBody(maxBytes int64) gin.HandlerFunc {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxRequestBodyBytes
	}
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
		default:
			c.Next()
			return
		}
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
