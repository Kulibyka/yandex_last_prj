package sqlite

import (
	"database/sql"
	"fmt"
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
