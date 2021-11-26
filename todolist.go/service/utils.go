package service

import (
	"fmt"
	s "sort"
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
	ftask.Limit = culcLimit(task.Deadline)
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

func formatTasksWithOption(tasks []database.Task, ctx *gin.Context) []FormatedTask {
	var ftasksWO []FormatedTask

	deadline, _ := ctx.GetQuery("deadline")
	sort, _ := ctx.GetQuery("sort")

	ftasks := formatTasks(tasks)
	if deadline == "yes" {
		for i:=0; i<len(ftasks); i++ {
			if ftasks[i].HasDeadline {
				ftasksWO = append(ftasksWO, ftasks[i])
			}
		}
	} else if deadline == "no" {
		for i:=0; i<len(ftasks); i++ {
			if !ftasks[i].HasDeadline {
				ftasksWO = append(ftasksWO, ftasks[i])
			}
		}
	} else {
		for i:=0; i<len(ftasks); i++ {
			ftasksWO = append(ftasksWO, ftasks[i])
		}
	}

	if sort == "reg_early" {
		s.SliceStable(ftasksWO, func(i, j int) bool { return ftasksWO[i].CreatedAt < ftasksWO[j].CreatedAt })
	} else if sort == "reg_late" {
		s.SliceStable(ftasksWO, func(i, j int) bool { return ftasksWO[i].CreatedAt > ftasksWO[j].CreatedAt })
	} else if sort == "dead_early" {
		s.SliceStable(ftasksWO, func(i, j int) bool { return ftasksWO[i].Deadline < ftasksWO[j].Deadline })
	} else if sort == "dead_late" {
		s.SliceStable(ftasksWO, func(i, j int) bool { return ftasksWO[i].Deadline > ftasksWO[j].Deadline })
	}

	return ftasksWO
}

func parseDeadline(deadline string) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	t, _ := time.ParseInLocation("2006-01-02T15:04", deadline, jst)
	return t
}

func culcLimit(deadline time.Time) string {
	sub := int(deadline.Sub(time.Now()).Seconds())
	negative := (sub<0)
	if negative {
		sub *= -1
	}

	seconds := sub%60
	sub -= seconds

	sub /= 60
	minutes := sub%60
	sub -= minutes

	sub /= 60
	hours := sub%24
	sub -= hours

	sub /= 24
	days := sub

	var res string
	if days >= 365 {
		res = fmt.Sprintf("%d年以上", int(days/365))
	} else if days >= 10 {
		res = fmt.Sprintf("%d日", days)
	} else {
		if days > 1 {
			res = fmt.Sprintf("%d日", days)
			if hours > 0 {
				res += fmt.Sprintf("%d時間", hours)
			}
		} else {
			if hours > 0 {
				res = fmt.Sprintf("%d時間", hours)
				res += fmt.Sprintf("%d分", minutes)
			} else {
				res = fmt.Sprintf("%d分", minutes)
			}
		}
	}

	if negative {
		return res + "前"
	} else {
		return res
	}
}
