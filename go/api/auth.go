package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aosh50/momenton/go/user"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Method %s, Path: %s", r.Method, r.URL.Path)
		if r.Method == "OPTIONS" {
			return
		}

		notAuthOnPost := []string{"/login", "/refresh"} //List of endpoints that don't require auth
		requestPath := r.URL.Path                       //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuthOnPost {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = Message("Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = Message("Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		token, tk, err := validateToken(tokenPart)
		if err != nil {
			response = Message(err.Error())
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
		}
		if !token.Valid {
			response = Message("Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			Respond(w, response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "Token", tk)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func validateToken(tokenString string) (*jwt.Token, *user.Token, error) {
	tk := &user.Token{}

	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	if err != nil {

		return nil, nil, errors.New("Malformed authentication token")
	}
	return token, tk, nil
}

func refresh(refreshToken string) (string, error) {
	token, tk, err := validateToken(refreshToken)
	if err != nil || !token.Valid {
		return "", errors.New("Token not valid")
	}
	return user.GenerateToken(time.Now().Add(24*time.Hour).Unix(), tk.User)
}
