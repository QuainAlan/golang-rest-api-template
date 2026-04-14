package auth

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "1234"
	hashedPassword, err := HashPassword(password)
	assert.Nil(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestGenerateToken(t *testing.T) {
	user := "chud"
	token, err := GenerateToken(user)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateTokenRoundTripUsernameClaim(t *testing.T) {
	const wantUser = "alice"
	tokenStr, err := GenerateToken(wantUser)
	if !assert.NoError(t, err) {
		return
	}

	parsed := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, parsed, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if !assert.NoError(t, err) {
		return
	}
	if !assert.True(t, token.Valid) {
		return
	}
	assert.Equal(t, wantUser, parsed.Username)
}

func TestGenerateRandomKey(t *testing.T) {
	randomKey := GenerateRandomKey()
	assert.NotEmpty(t, randomKey)
	assert.Len(t, randomKey, 44)
}
