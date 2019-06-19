package shared

import (
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// JWT is the encryption key used to sign the JWTs
// TODO: It must be configurable via env variable
var JWT = []byte("super_secret_key")

// Claims represents the jwt decoded
type Claims struct {
	UUID string `json:"uuid"`
	jwt.StandardClaims
}

// HashPassword hashes the password with bcrypt
func HashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		GetLogger().Error("Could not hash user's password", zap.Error(err))
	}
	return string(hash), err
}

// CheckPassword checks if the provided passwd
// could match the hashed password
func CheckPassword(hashed, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(passwd))
	if err != nil {
		return false
	}

	return true
}
