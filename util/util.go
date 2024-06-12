package util

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"

	"github.com/Terracode-Dev/terraui-back/types"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSecretKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func UKGen(s string) string {
	return "uk12iu12"
}

func UKcheck(id string, time int64, key string) bool {
	return true
}

func HashKey(key string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	return string(bytes), err
}

func HashCheck(key, pkey string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pkey), []byte(key))
	return err == nil
}

func GetToken(user *types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": user.Userid,
		"role":   user.Userrole,
		"uk":     "aa11bb22",
		"email":  user.Useremail,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JKEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
