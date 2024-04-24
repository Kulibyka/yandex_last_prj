package getOperation

import (
	"encoding/json"
	"gRPC_service/internal/calculator"
	"net/http"
)

func GetOperations(w http.ResponseWriter, r *http.Request) {
	// Создаем словарь для хранения имени операции и ее длительности в секундах
	operationDurations := make(map[string]int64)

	// Проходим по каждой операции
	for _, op := range calculator.Operations {
		duration := op.Duration.Seconds()
		operationDurations[op.Name] = int64(duration)
	}

	json.NewEncoder(w).Encode(operationDurations)
}
