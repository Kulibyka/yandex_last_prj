package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"gRPC_service/internal/handlers/addTask"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS Users (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        Login TEXT NOT NULL,
        Password TEXT NOT NULL
    );`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`
    CREATE TABLE IF NOT EXISTS Tasks (
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        UserID INTEGER NOT NULL,
        Expression TEXT NOT NULL,
        Status TEXT NOT NULL,
        Result REAL,
        AgentID INTEGER,
        Duration REAL,
		StartDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	EndDate TIMESTAMP,
        FOREIGN KEY (UserID) REFERENCES Users(ID)
	);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Register(login string, password string) (int64, error) {
	const op = "storage.sqlite.Register"

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM Users WHERE Login=?", login).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if count > 0 {
		return 0, fmt.Errorf("%s: user already exist", op)
	}

	res, err := s.db.Exec("INSERT INTO Users (Login, Password) VALUES (?, ?)", login, password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) Login(login string, password string) (int64, error) {
	const op = "storage.sqlite.Login"

	var userID int64
	var dbPassword string
	err := s.db.QueryRow("SELECT ID, Password FROM Users WHERE Login=?", login).Scan(&userID, &dbPassword)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, совпадает ли пароль
	if password != dbPassword {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) SaveTask(userID int64, expression string) (int64, error) {
	const op = "storage.sqlite.SaveTask"
	stmt, err := s.db.Prepare("INSERT INTO Tasks (UserID, Expression, Status) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(userID, expression, "queued")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) GetTaskFromDB(agentID int64) (*addTask.Task, error) {
	const op = "storage.sqlite.GetTaskFromDB"
	var task addTask.Task
	// Выполняем запрос к базе данных для получения задачи со статусом "queued"
	row := s.db.QueryRow("SELECT ID, Expression FROM Tasks WHERE Status = 'queued' LIMIT 1")

	// Сканируем данные из строки результата в переменные
	if err := row.Scan(&task.ID, task.Expression); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Если нет ни одной задачи со статусом "queued", возвращаем nil
			return nil, nil
		}
		// В случае других ошибок возвращаем их
		return nil, err
	}

	_, err := s.db.Exec("UPDATE Tasks SET AgentID=?, Status = 'doing calculations', StartDate = ? WHERE ID = ?",
		agentID, time.Now(), task.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &task, nil
}

func (s *Storage) UpdateTaskInDB(task *addTask.Task) error {
	const op = "storage.sqlite.UpdateTaskInDB"
	// Выполняем запрос к базе данных для обновления задачи
	_, err := s.db.Exec("UPDATE Tasks SET Status=?, Result=?, Duration=?, AgentID=?, EndDate=? WHERE ID=?",
		task.Status, task.Result, task.Duration, task.AgentID, time.Now(), task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetResult(taskID int) (*addTask.Task, error) {
	const op = "storage.sqlite.GetResult"
	var task addTask.Task

	err := s.db.QueryRow("SELECT Status, Result FROM Tasks WHERE ID = ?",
		taskID).Scan(&task.Status, &task.Result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &task, nil
}
