package register

import (
	"encoding/json"
	resp "gRPC_service/lib/response"
	"github.com/go-chi/render"
	"net/http"
)

type UserRegister interface {
	Register(login string, password string) (int64, error)
}

func NewUser(userRegister UserRegister) http.HandlerFunc {
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

		// Проверяем, что логин и пароль не пустые
		if user.Login == "" || user.Password == "" {
			http.Error(w, "Login and password cannot be empty", http.StatusBadRequest)
			return
		}

		_, err := userRegister.Register(user.Login, user.Password)
		if err != nil {
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
	}
}
