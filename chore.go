package main

import "time"

type Chore struct {
	ID          int       `json:"id"`
	Created     time.Time `json:"created"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	Done        bool      `json:"done"`
	Execution   int       `json:"execution"`
}

type Execution struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	ChoreID  int    `json:"chore_id"`
	Points   int    `json:"points"`
	UserID   int    `json:"done_by"`
}
