package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	Deadline  time.Time `db:"deadline"`
	IsDone    bool      `db:"is_done"`
	Detail    string    `db:"detail"`
}

// User corresponds to a row in `users` table
type User struct {
	ID       uint64 `db:"id"`
	Username string `db:"username"`
	Passward string `db:"passward"`
}

// User corresponds to a row in `owners` table
type Owner struct {
	ID       uint64 `db:"id"`
	Username string `db:"username"`
	TaskID   uint64 `db:"taskid"`
}
