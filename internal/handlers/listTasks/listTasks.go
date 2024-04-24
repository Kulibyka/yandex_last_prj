package listTasks

import (
	"encoding/json"
	"gRPC_service/internal/handlers/addTask"
	"net/http"
)

func ListTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(addTask.Tasks)
}
