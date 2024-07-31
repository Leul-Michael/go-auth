package utils

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// If phone number is empty, consider it as valid (no validation required)
	if phoneNumber == "" {
		return true
	}
	// Adjust the regex pattern to match ethiopian phone number format
	regex := regexp.MustCompile(`^(?:\+251|0)[79]\d{8}$`)
	return regex.MatchString(phoneNumber)
}

func GenerateToken(duration int, sub interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour * time.Duration(duration)).Unix(),
	})

	generatedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}
	return generatedToken, nil
}
