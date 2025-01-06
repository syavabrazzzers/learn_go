package utils

import (
	"fmt"
	"learn/models"
	"math/rand/v2"

	settings "learn/settings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func MakePasswordHash(pass string) string {
	password_hash, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		panic(err)
	}
	return string(password_hash)
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MakeJwt(user models.User) string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": fmt.Sprint(user.ID),
		"iss": "learn-app",
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte(settings.Settings.JwtSecret))
	if err != nil {
		panic(err)
	}

	return token
}

func VerifyJwt(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.Settings.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return parsedToken, nil
}

func MakeVerificationCode() int {
	return rand.IntN(9999-1000) + 1000
}
