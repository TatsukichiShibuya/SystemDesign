package service

import (
	"fmt"
	"time"
	"crypto/sha256"

	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func sessionCheck(ctx *gin.Context) bool {
  session := sessions.Default(ctx)
  return session.Get("username") != nil
}

func hash(passward string) string {
	h := sha256.Sum256([]byte(passward))
	return fmt.Sprintf("%x", h)
}

type FormatedTask struct {
	ID             uint64
	Title          string
	CreatedAt      string
	CreatedAt_html string
	Deadline       string
	Deadline_html  string
	HasDeadline    bool
	Limit				   string
	IsDone         bool
	Detail         string
}

func formatTask(task database.Task) FormatedTask {
	var ftask FormatedTask

	ftask.ID = task.ID
	ftask.Title = task.Title
	ftask.CreatedAt = task.CreatedAt.Format("2006/01/02 15:04")
	ftask.CreatedAt_html = task.CreatedAt.Format("2006-01-02T15:04")
	ftask.Deadline = task.Deadline.Format("2006/01/02 15:04")
	ftask.Deadline_html = task.Deadline.Format("2006-01-02T15:04")
	if task.Deadline.After(task.CreatedAt) {
		ftask.HasDeadline = true
	} else {
		ftask.HasDeadline = false
	}
	ftask.Limit = "-"
	ftask.IsDone = task.IsDone
	ftask.Detail = task.Detail

	return ftask
}

func formatTasks(tasks []database.Task) []FormatedTask {
	ftasks := make([]FormatedTask, len(tasks))

	for i:=0; i<len(tasks); i++ {
		ftasks[i] = formatTask(tasks[i])
	}

	return ftasks
}

func parseDeadline(deadline string) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	t, _ := time.ParseInLocation("2006-01-02T15:04", deadline, jst)
	return t
}
