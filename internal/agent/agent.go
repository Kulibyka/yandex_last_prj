package agent

import (
	"gRPC_service/internal/calculator"
	"gRPC_service/internal/handlers/addTask"
	"time"
)

func StartAgents(numAgents int) {
	for i := 0; i < numAgents; i++ {
		go func(agentID int) {
			for task := range addTask.TaskChannel {
				task.AgentID = agentID

				task.Status = "doing calculations"
				startTime := time.Now()

				task.Mutex.Lock()
				task.Result = calculator.EvaluateExpression(task.Expression)
				task.Mutex.Unlock()

				task.Status = "completed"
				endTime := time.Now()
				task.Duration = endTime.Sub(startTime).Seconds()
			}

		}(i + 1)
	}
}
