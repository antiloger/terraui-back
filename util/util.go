package util

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateSecretKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func GetToken(userID string, tenID string, Skey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"tid": tenID,
		"exp": time.Now().Add(time.Hour * 24).Unix()})
	tokenString, err := token.SignedString([]byte(Skey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
