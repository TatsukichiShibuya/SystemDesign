package service

import (
	"fmt"
	s "sort"
	"time"
	"strings"
	"unicode"
	"crypto/sha256"

	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
	"github.com/jmoiron/sqlx"

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

func checkname(username string) bool {
	res := true
	for _, c := range username {
		ok := false
		ok = ok || unicode.In(c, unicode.Hiragana)
		ok = ok || unicode.In(c, unicode.Katakana)
		ok = ok || unicode.In(c, unicode.Han)
		ok = ok || unicode.IsDigit(c)
		ok = ok || unicode.IsLetter(c)
		res = res && ok
	}
	return res
}

func parseUsers(users string) []string {
	var res []string
	if users == "" {
		return res
	} else {
		return strings.Split(users, ",")
	}
}

func allExist(users []string, db *sqlx.DB) bool {
	var user database.User
	for _, username := range users {
		if checkname(username) {
			err := db.Get(&user, "SELECT * FROM users WHERE username=?", username)
			if err != nil {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

type FormatedTask struct {
	ID             uint64
	Title          string
	CreatedAt      string
	CreatedAt_html string
	Deadline       string
	Deadline_html  string
	HasDeadline    bool
	Importance     int
	Limit				   string
	IsShared       bool
	SharedUsers    string
	IsDone         bool
	Detail         string
}

func formatTask(task database.Task, ctx *gin.Context, db *sqlx.DB, username string) (FormatedTask, error) {
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
	if ftask.HasDeadline && !ftask.IsDone {
		ftask.Importance = culcImportance(task.Deadline)
	} else {
		ftask.Importance = 0
	}

	var owners []database.Owner
	err := db.Select(&owners, "SELECT * FROM owners WHERE taskid=?", task.ID)
	if err != nil {
		return ftask, err
	}
	ftask.IsShared = (len(owners)>1)
	ftask.SharedUsers = ""
	for _, owner := range owners {
		if owner.Username != username {
			if ftask.SharedUsers != "" {
				ftask.SharedUsers += ","
			}
			ftask.SharedUsers += owner.Username
		}
	}


	ftask.Detail = task.Detail

	return ftask, nil
}

func formatTasks(tasks []database.Task, ctx *gin.Context) ([]FormatedTask, error) {
	ftasks := make([]FormatedTask, len(tasks))

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		return ftasks, err
	}

	session := sessions.Default(ctx)
	username := session.Get("username").(string)

	for i:=0; i<len(tasks); i++ {
		ftasks[i], err = formatTask(tasks[i], ctx, db, username)
		if err != nil {
			return ftasks, err
		}
	}

	return ftasks, nil
}

func formatTasksWithOption(tasks []database.Task, ctx *gin.Context) ([]FormatedTask, error) {
	var ftasksWO []FormatedTask

	deadline, _ := ctx.GetQuery("deadline")
	sort, _ := ctx.GetQuery("sort")

	ftasks, err := formatTasks(tasks, ctx)
	if err != nil {
		return ftasksWO, err
	}

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
		s.SliceStable(ftasksWO, func(i, j int) bool {
			if (ftasksWO[i].HasDeadline&&ftasksWO[j].HasDeadline) || (!ftasksWO[i].HasDeadline&&!ftasksWO[j].HasDeadline){
				return ftasksWO[i].Deadline < ftasksWO[j].Deadline
			} else {
				return ftasksWO[i].HasDeadline
			}
		})
	} else if sort == "dead_late" {
		s.SliceStable(ftasksWO, func(i, j int) bool {
			if (ftasksWO[i].HasDeadline&&ftasksWO[j].HasDeadline) || (!ftasksWO[i].HasDeadline&&!ftasksWO[j].HasDeadline){
				return ftasksWO[i].Deadline > ftasksWO[j].Deadline
			} else {
				return ftasksWO[i].HasDeadline
			}
		})
	}

	return ftasksWO, nil
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


func culcImportance(deadline time.Time) int {
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

	if negative {
		return 3
	} else if days == 0 {
		return 2
	} else {
		return 1
	}
}
