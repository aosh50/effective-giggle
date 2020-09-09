package user

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	jwt.StandardClaims
	User string
}

func dbPassword() string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("password")), bcrypt.DefaultCost)
	if err != nil {
		logrus.Fatal(err)
	}
	return string(hashedPassword)
}

func Login(user, password string) (string, string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword()), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return "", "", errors.New("Invalid login credentials. Please try again")
	}

	accessToken, err := GenerateToken(time.Now().Add(24*time.Hour).Unix(), user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenerateToken(time.Now().AddDate(0, 0, 7).Unix(), user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil

}

func GenerateToken(exp int64, user string) (string, error) {
	tk := &Token{
		jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    "code-test",
		},
		user,
	}

	tok := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	return tok.SignedString([]byte(os.Getenv("token_password")))

}
