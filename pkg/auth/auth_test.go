package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	token, err := jwt.ParseWithClaims(tokenStr, parsed, JWTKeyFunc(JwtKey))
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

func TestJWTKeyFuncRejectsNonHS256Algorithms(t *testing.T) {
	key := []byte("jwt-keyfunc-test-secret-32bytes!!")
	exp := time.Now().Add(time.Hour)
	baseClaims := Claims{
		Username: "attacker",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t.Run("HS512", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS512, &baseClaims)
		s, err := tok.SignedString(key)
		if !assert.NoError(t, err) {
			return
		}
		parsed := &Claims{}
		_, err = jwt.ParseWithClaims(s, parsed, JWTKeyFunc(key))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})

	t.Run("none", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodNone, &baseClaims)
		s, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if !assert.NoError(t, err) {
			return
		}
		parsed := &Claims{}
		_, err = jwt.ParseWithClaims(s, parsed, JWTKeyFunc(key))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})
}

func TestJWTKeyFuncAcceptsHS256(t *testing.T) {
	key := []byte("jwt-keyfunc-hs256-accept-secret!")
	claims := &Claims{
		Username: "legit",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(key)
	if !assert.NoError(t, err) {
		return
	}
	parsed := &Claims{}
	token, err := jwt.ParseWithClaims(s, parsed, JWTKeyFunc(key))
	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, token.Valid)
	assert.Equal(t, "legit", parsed.Username)
}
