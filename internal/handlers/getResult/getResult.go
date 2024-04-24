package getResult

import (
	"encoding/json"
	"gRPC_service/internal/handlers/addTask"
	"net/http"
	"strconv"
	"strings"
)

func GetTaskResult(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	taskID, err := strconv.Atoi(segments[0])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if taskID < 1 || taskID > len(addTask.Tasks) {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Получаем задачу по ее ID
	task := addTask.Tasks[taskID-1]

	// Если задача еще не завершена, отправляем сообщение ожидания
	if task.Status != "completed" {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "Task is not completed yet"})
		return
	}

	response := map[string]float64{"result": task.Result}
	json.NewEncoder(w).Encode(response)
}
