package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

// Newtask renders a task with given ID
func GetNewtask(ctx *gin.Context) {
	if sessionCheck(ctx) {
		Newtask(ctx, "")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostNewtask(ctx *gin.Context) {
	if sessionCheck(ctx) {
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		title, _ := ctx.GetPostForm("title")
		detail, _ := ctx.GetPostForm("detail")
		message := ""

		if title == "" {
			message += "タスク名を記入してください"
		} else {
			// register new task
			data := map[string]interface{}{ "title": title, "detail": detail }
			res, _ := db.NamedExec("INSERT INTO tasks (title, detail) VALUES (:title, :detail);", data)
			// set owner
			session := sessions.Default(ctx)
			username := session.Get("username")
			lastid, _ := res.LastInsertId()
			data = map[string]interface{}{ "username": username, "taskid": lastid }
			_, _ = db.NamedExec("INSERT INTO owners (username, taskid) VALUES (:username, :taskid);", data)
		}

		if message == "" {
			ctx.Redirect(http.StatusSeeOther, "/list")
		} else {
			Newtask(ctx, message)
		}
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func Newtask(ctx *gin.Context, message string) {
	ctx.HTML(http.StatusOK, "newtask.html", gin.H{ "Title": "NEW TASK",
																					       "Message": message })
}
