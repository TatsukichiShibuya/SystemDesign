package service

import (
	s "sort"

	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
	"github.com/jmoiron/sqlx"

	database "todolist.go/db"
)

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

const FORMAT_LIST = "2006-01-02T15:04"
const FORMAT_HTML = "2006/01/02 15:04"

func formatTask(task database.Task, db *sqlx.DB, userid uint64) (FormatedTask, error) {
	var ftask FormatedTask

	ftask.ID = task.ID
	ftask.Title = task.Title
	ftask.CreatedAt = task.CreatedAt.Format(FORMAT_LIST)
	ftask.CreatedAt_html = task.CreatedAt.Format(FORMAT_HTML)
	ftask.Deadline = task.Deadline.Format(FORMAT_LIST)
	ftask.Deadline_html = task.Deadline.Format(FORMAT_HTML)
	ftask.HasDeadline = task.Deadline.After(task.CreatedAt)
	ftask.Limit = culcLimit(task.Deadline)
	ftask.IsDone = task.IsDone
	if ftask.HasDeadline && !ftask.IsDone {
		ftask.Importance = culcImportance(task.Deadline)
	} else {
		ftask.Importance = 0
	}
  ftask.Detail = task.Detail

	var owners []database.Owner
	err := db.Select(&owners, "SELECT * FROM owners WHERE taskid=?", task.ID)
	if err != nil {
		return ftask, err
	}
	ftask.IsShared = len(owners)>=2
  ftask.SharedUsers = ""
  var suser database.User
	for _, owner := range owners {
		if owner.UserID != userid {
      // get shareuser's name
			err = db.Get(&suser, "SELECT * FROM users WHERE id=?", owner.UserID)
			if err != nil {
				return ftask, err
			}

      if ftask.SharedUsers != "" {
				ftask.SharedUsers += ","
			}
			ftask.SharedUsers += suser.Username
		}
	}

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
	userid := session.Get("userid").(uint64)

  // format tasks
	for i:=0; i<len(tasks); i++ {
		ftasks[i], err = formatTask(tasks[i], db, userid)
		if err != nil {
			return ftasks, err
		}
	}

	return ftasks, nil
}

func formatTasksWithOption(tasks []database.Task, ctx *gin.Context) ([]FormatedTask, error) {
	var ftasksWO []FormatedTask

	sort, _ := ctx.GetQuery("sort")

	ftasks, err := formatTasks(tasks, ctx)
	if err != nil {
		return ftasksWO, err
	}

	for i:=0; i<len(ftasks); i++ {
		ftasksWO = append(ftasksWO, ftasks[i])
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
