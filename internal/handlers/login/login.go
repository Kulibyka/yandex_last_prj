package login

import (
	"encoding/json"
	"gRPC_service/internal/auth"
	resp "gRPC_service/lib/response"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type UserLogin interface {
	Login(login string, password string) (int64, error)
}

func Login(userLogin UserLogin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if user.Login == "" || user.Password == "" {
			http.Error(w, "Login and password cannot be empty", http.StatusBadRequest)
			return
		}

		id, err := userLogin.Login(user.Login, user.Password)
		if err != nil {
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		token, err := auth.CreateToken(id)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// Отправляем токен в ответе
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
