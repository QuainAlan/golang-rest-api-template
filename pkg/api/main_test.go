package api

import (
	"bytes"
	"os"
	"testing"

	"golang-rest-api-template/pkg/auth"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	if err := auth.SetJWTSigningKey(bytes.Repeat([]byte("t"), auth.MinJWTSecretKeyBytes)); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
