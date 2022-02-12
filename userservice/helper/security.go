package helper

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserDetails struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthClaim struct {
	jwt.StandardClaims
	UserDetails
}

func NewUserDetails(uuid string, name string, email string) UserDetails {
	return UserDetails{
		Uuid:  uuid,
		Name:  name,
		Email: email,
	}
}

func HashPassword(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ValidatePassword(reqPasswd string, orgPasswdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(orgPasswdHash), []byte(reqPasswd))

	return err == nil
}

func CreateJwt(userDetails UserDetails) (string, error) {
	claims := AuthClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    GetEnv("JWT_ISSUER"),
		},
		userDetails,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(GetEnv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJwt(tokenString string) (userDetails UserDetails, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*AuthClaim); ok && token.Valid {
		return claims.UserDetails, nil
	}

	return userDetails, errors.New("Token Expired.")
}

func GenerateUuid() string {
	return uuid.New().String()
}
