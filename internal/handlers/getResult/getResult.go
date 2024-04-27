package getResult

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gRPC_service/internal/handlers/addTask"
	"net/http"
	"strconv"
	"strings"
)

type ResultGetter interface {
	GetResult(int) (*addTask.Task, error)
}

func GetTaskResult(resultGetter ResultGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
		taskID, err := strconv.Atoi(segments[0])
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		task, err := resultGetter.GetResult(taskID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}
			fmt.Println(err)
			http.Error(w, "Failed to fetch task details", http.StatusInternalServerError)
			return
		}

		// Если задача еще не завершена, отправляем сообщение ожидания
		if task.Status != "completed" {
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]string{"status": "Task is not completed yet"})
			return
		}

		response := map[string]float64{"result": task.Result}
		json.NewEncoder(w).Encode(response)
	}
}
