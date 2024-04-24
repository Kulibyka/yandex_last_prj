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
)

var (
	Tasks        []*Task
	nextTaskID   = 1
	TaskChannel  = make(chan *Task, 100)
	mySigningKey = []byte("yandexL")
)

type Task struct {
	ID         int        `json:"id"`
	Expression string     `json:"expression"`
	Status     string     `json:"status"`
	Result     float64    `json:"result,omitempty"`
	AgentID    int        `json:"agent_id,omitempty"`
	Duration   float64    `json:"duration(seconds),omitempty"`
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
		defer r.Body.Close()
		// Добавляем задачу в список
		task.ID = nextTaskID
		nextTaskID++
		task.Status = "queued"
		task.Mutex = sync.Mutex{}
		Tasks = append(Tasks, &task)

		userID, err := getUserID(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))

		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusForbidden)
			log.Println(err)
			return
		}

		_, err = taskSaver.SaveTask(userID, task.Expression)

		TaskChannel <- &task

		json.NewEncoder(w).Encode(map[string]int{"task_id": task.ID})
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
