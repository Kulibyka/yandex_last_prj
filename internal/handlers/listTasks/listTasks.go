package listTasks

import (
	"encoding/json"
	"gRPC_service/internal/handlers/addTask"
	"net/http"
)

type ListGetter interface {
	GetTaskList() ([]*addTask.Task, error)
}

func ListTasks(listGetter ListGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := listGetter.GetTaskList()
		if err != nil {
			http.Error(w, "Failed to get task list", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(tasks)
	}
}
