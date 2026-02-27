package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDGenerated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		id1 := GetRequestID(c)
		id2 := GetRequestID(c)
		ctxID := GetRequestIDFromContext(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{"id1": id1, "id2": id2, "ctx": ctxID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	requestID := w.Header().Get(RequestIDHeader)
	assert.NotEmpty(t, requestID)

	var body map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, requestID, body["id1"])
	assert.Equal(t, requestID, body["id2"])
	assert.Equal(t, requestID, body["ctx"])
}

func TestRequestIDPreserved(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		requestID := GetRequestID(c)
		ctxID := GetRequestIDFromContext(c.Request.Context())
		c.JSON(http.StatusOK, gin.H{"id": requestID, "ctx": ctxID})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(RequestIDHeader, "req-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "req-123", w.Header().Get(RequestIDHeader))

	var body map[string]string
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "req-123", body["id"])
	assert.Equal(t, "req-123", body["ctx"])
}
