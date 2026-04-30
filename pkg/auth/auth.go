package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Claims struct to be encoded to JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// JWTKeyFunc returns a jwt.Keyfunc for ParseWithClaims that only accepts
// tokens signed with HMAC-SHA256 using key. Any other signing method
// (including "none" and asymmetric algorithms) is rejected to prevent
// algorithm confusion attacks.
func JWTKeyFunc(key []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateToken(username string) (string, error) {
	expiresAt := time.Now().Add(5 * time.Minute).Unix()
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRandomKey() string {
	key := make([]byte, 32) // generate a 256 bit key
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}
