package utils

import (
	"fmt"
	"learn/db/redis"
	"learn/models"
	"learn/utils/smtp"
	"math/rand/v2"
	"net/http"

	settings "learn/settings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

// func MakeVerificationCode() int {
// 	return rand.IntN(9999-1000) + 1000
// }

func GetRequestBody[T any](ctx *gin.Context) (*T, error) {
	var data T
	if err := ctx.BindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"detail": "Invalid request body"})
		return nil, err
	}
	return &data, nil
}

func SendVerificationCode(email string) (verification_key string, verification_code string) {
	verification_key = uuid.New().String()
	verification_code = fmt.Sprint(rand.IntN(9999-1000) + 1000)

	go redis.Client.SetJson(
		verification_key,
		map[string]string{
			"email": email,
			"code":  verification_code,
		},
		settings.Settings.VerificationCodeExpiration,
	)
	go smtp.SendMail([]string{email}, verification_code, "")
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%&.?")

func GenerateRecoveryCodes() []string {
	result := make([]string, settings.Settings.RecoveryCodeCount)

	for i := range settings.Settings.RecoveryCodeCount {
		b := make([]rune, settings.Settings.RecoveryCodeLength)
		for j := range b {
			b[j] = letters[rand.IntN(len(letters))]
		}
		result[i] = string(b)
	}

	return result
}
