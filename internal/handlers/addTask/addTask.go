package addTask

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	Tasks        []*Task
	mySigningKey = []byte("yandexL")
)

type Task struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Expression string     `json:"expression"`
	Status     string     `json:"status"`
	Result     float64    `json:"result,omitempty"`
	AgentID    int        `json:"agent_id,omitempty"`
	Duration   float64    `json:"duration(seconds),omitempty"`
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date,omitempty"`
	Mutex      sync.Mutex `json:"-"`
}

type TaskSaver interface {
	SaveTask(userID int64, expression string) (int64, error)
}

func AddTask(taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}
		// Добавляем задачу в список

		userID, err := getUserID(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusForbidden)
			log.Println(err)
			return
		}

		id, err := taskSaver.SaveTask(userID, task.Expression)

		json.NewEncoder(w).Encode(map[string]int64{"task_id": id})
	}
}

func getUserID(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return mySigningKey, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid token claims")
	}

	userID, ok := claims["userID"].(float64)
	if !ok {
		return 0, errors.New("User ID not found in token claims")
	}

	return int64(userID), nil
}
