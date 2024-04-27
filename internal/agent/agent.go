package agent

import (
	"gRPC_service/internal/calculator"
	"gRPC_service/internal/handlers/addTask"
	"log"
	"time"
)

type DbChanger interface {
	GetTaskFromDB(agentID int64) (*addTask.Task, error)
	UpdateTaskInDB(task *addTask.Task) error
}

func StartAgents(numAgents int, dbChanger DbChanger) {
	for i := 0; i < numAgents; i++ {
		go func(agentID int) {
			for {
				task, err := dbChanger.GetTaskFromDB(int64(agentID))
				if err != nil {
					log.Printf("Error getting task from database: %v", err)
					continue
				}

				// Проверяем, есть ли задача
				if task == nil {
					// Если задачи нет, ждем некоторое время перед повторной попыткой
					time.Sleep(5 * time.Second)
					continue
				}

				task.Status = "doing calculations"
				startTime := time.Now()

				task.Mutex.Lock()
				task.Result = calculator.EvaluateExpression(task.Expression)
				task.Mutex.Unlock()

				task.Status = "completed"
				endTime := time.Now()
				task.EndDate = endTime
				task.Duration = endTime.Sub(startTime).Seconds()

				err = dbChanger.UpdateTaskInDB(task)
				if err != nil {
					log.Printf("Error updating task in database: %v", err)
					continue
				}
			}
		}(i + 1)
	}
}
