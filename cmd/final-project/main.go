package main

import (
	"gRPC_service/internal/agent"
	"gRPC_service/internal/auth"
	"gRPC_service/internal/handlers/addTask"
	"gRPC_service/internal/handlers/getOperation"
	"gRPC_service/internal/handlers/getResult"
	"gRPC_service/internal/handlers/listTasks"
	"gRPC_service/internal/handlers/login"
	"gRPC_service/internal/handlers/register"
	"gRPC_service/internal/storage/sqlite"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {

	storage, err := sqlite.New("./storage/storage.db")
	if err != nil {
		log.Fatal("failed to init storage: %w", err)
	}

	go agent.StartAgents(3, storage)

	router := mux.NewRouter()
	router.HandleFunc("/register", register.NewUser(storage)).Methods("POST")
	router.HandleFunc("/login", login.Login(storage)).Methods("POST")

	router.Handle("/tasks/add", auth.IsAuthorized(addTask.AddTask(storage))).Methods("POST")
	router.Handle("/tasks", auth.IsAuthorized(listTasks.ListTasks)).Methods("GET")
	router.Handle("/tasks/{id}/result", auth.IsAuthorized(getResult.GetTaskResult(storage))).Methods("GET")
	router.Handle("/operations", auth.IsAuthorized(getOperation.GetOperations)).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
