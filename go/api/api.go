package api

import (
	"encoding/json"
	"net/http"
	"os"

	user "github.com/aosh50/momenton/go/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	godotenv.Load() //Load .env file

}
func Start() {
	router := Router()
	port := os.Getenv("PORT")

	if port == "" {
		port = "3333"
	}
	logrus.Infof("Listening on port %s", port)
	http.ListenAndServe(":"+port, router)
}

func Router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(JwtAuthentication)

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		type LoginModel struct {
			User     string
			Password string
		}
		login := &LoginModel{}
		err := json.NewDecoder(r.Body).Decode(&login) //decode the request body into struct
		if err != nil {
			Respond(w, Message("Invalid request"))
			logrus.Error(err.Error())
			return
		}
		token, refreshToken, err := user.Login(login.User, login.Password)
		if err != nil {
			Respond(w, Message(err.Error()))
			logrus.Error(err.Error())
			return
		}
		data := map[string]interface{}{
			"token":         token,
			"refresh_token": refreshToken,
		}

		Respond(w, data)

	})
	r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value("Token").(*user.Token)
		Respond(w, Message(token.User))

	})
	r.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		type RefreshModel struct {
			RefreshToken string `json:"refresh_token"`
		}
		model := &RefreshModel{}
		err := json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			Respond(w, Message("Invalid request"))
			logrus.Error(err.Error())
			return
		}
		logrus.Info(model.RefreshToken)
		token, err := refresh(model.RefreshToken)
		if err != nil {
			Respond(w, Message(err.Error()))
			logrus.Error(err.Error())
			return
		}
		data := map[string]interface{}{
			"token": token,
		}

		Respond(w, data)

	})
	return r
}
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if val, ok := data["status"]; ok {
		if !val.(bool) {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	json.NewEncoder(w).Encode(data)
}
func Message(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}
