package main

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type ChoreStorage interface {
	// List all items
	List() ([]Chore, error)
	ListTodo() ([]Chore, error)
	// Get a single item
	Get(id int) (Chore, error)
	// Create a new item
	Create(item Chore) (int, error)
	// Update an existing item
	Update(id int, item Chore) error
	// Complete an item
	Complete(id int, execution *Execution) error
	// Delete an item
	Delete(id int) error
}

type SqlLiteStorage struct {
	db *sql.DB
}

const createChores string = `
  CREATE TABLE IF NOT EXISTS chores (
  id INTEGER NOT NULL PRIMARY KEY,
  created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name TEXT,
  description TEXT,
	points INTEGER DEFAULT 1,
	done BOOLEAN DEFAULT FALSE,
	execution INTEGER DEFAULT 0
  );`

const createExecution string = `
  CREATE TABLE IF NOT EXISTS execution (
  id INTEGER NOT NULL PRIMARY KEY,
	done_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	name TEXT,
	user_name TEXT,
	chore_id	INTEGER DEFAULT 0,
	points INTEGER DEFAULT 1,
  user_id INTEGER DEFAULT 0
  );`

const createUsers string = `
	CREATE TABLE IF NOT EXISTS user (
  id INTEGER NOT NULL PRIMARY KEY,
	name TEXT
  );`

// ID       int    `json:"id"`
// Name     string `json:"name"`
// UserName string `json:"user_name"`
// ChoreID  int    `json:"chore_id"`
// Points   int    `json:"points"`
// UserID   int    `json:"done_by"`

func NewSqlStorage() (*SqlLiteStorage, error) {
	db, err := sql.Open("sqlite3", filepath.Join("data", "chores.db"))
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createChores)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createExecution)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(createUsers)
	if err != nil {
		return nil, err
	}
	return &SqlLiteStorage{
		db: db,
	}, nil
}

func (s *SqlLiteStorage) List() ([]Chore, error) {
	res, err := s.db.Query("SELECT id, created, name, description, points, done, execution FROM chores")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var items []Chore
	for res.Next() {
		var item Chore
		if err := res.Scan(&item.ID, &item.Created, &item.Name, &item.Description, &item.Points, &item.Done, &item.Execution); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil

}

func (s *SqlLiteStorage) ListTodo() ([]Chore, error) {
	res, err := s.db.Query("SELECT id, created, name, description, points, done, execution FROM chores WHERE done = false")
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var items []Chore
	for res.Next() {
		var item Chore
		if err := res.Scan(&item.ID, &item.Created, &item.Name, &item.Description, &item.Points, &item.Done, &item.Execution); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil

}

func (s *SqlLiteStorage) Get(id int) (Chore, error) {
	res := s.db.QueryRow("SELECT id, created, name, description, points, done, execution FROM chores WHERE id = ?", id)

	var item Chore
	err := res.Scan(&item.ID, &item.Created, &item.Name, &item.Description, &item.Points, &item.Done, &item.Execution)

	return item, err
}

func (s *SqlLiteStorage) Create(item Chore) (int, error) {
	res, err := s.db.Exec("INSERT INTO chores (name, description, points) VALUES (?, ?, ?)", item.Name, item.Description, item.Points)
	if err != nil {
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SqlLiteStorage) Update(id int, item Chore) error {
	_, err := s.db.Exec("UPDATE chores SET name = ?, description = ?, points = ? WHERE id = ?", item.Name, item.Description, item.Points, id)
	return err
}

func (s *SqlLiteStorage) Complete(id int, execution *Execution) error {
	item, err := s.Get(id)
	if err != nil {
		return err
	}

	if item.Done {
		return nil
	}

	res, err := s.db.Exec("INSERT INTO execution (name, userName, chore_id, points, user_id) VALUES (?, ?, ?, ?, ?)", execution.Name, execution.UserName, id, item.Points, execution.UserID)
	if err != nil {
		return err
	}
	var eid int64
	if eid, err = res.LastInsertId(); err != nil {
		return err
	}
	_, err = s.db.Exec("UPDATE chores SET done = ?, execution = ? WHERE id = ?", true, eid, id)

	return err
}

func (s *SqlLiteStorage) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM chores WHERE id = ?", id)
	return err
}
