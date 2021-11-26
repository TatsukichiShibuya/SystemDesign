package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetNewtask(ctx *gin.Context) {
	if sessionCheck(ctx) {
		Newtask(ctx, "")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostNewtask(ctx *gin.Context) {
	if sessionCheck(ctx) {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// get inputs
		title, _ := ctx.GetPostForm("title")
		deadline_string, _ := ctx.GetPostForm("deadline")
		detail, _ := ctx.GetPostForm("detail")
		message := ""

		if title == "" {
			message = "タスク名を記入してください"
		} else {
			// parse deadline
			deadline := parseDeadline(deadline_string)

			// register new task
			data := map[string]interface{}{ "title": title, "deadline": deadline, "detail": detail}
			var query string
			if deadline.IsZero() {
				query = "INSERT INTO tasks (title, detail) VALUES (:title, :detail);"
			} else {
				query = "INSERT INTO tasks (title, deadline, detail) VALUES (:title, :deadline, :detail);"
			}
			res, err := db.NamedExec(query, data)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			// set owner
			session := sessions.Default(ctx)
			username := session.Get("username")
			lastid, err := res.LastInsertId()
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{ "username": username, "taskid": lastid }
			_, err = db.NamedExec("INSERT INTO owners (username, taskid) VALUES (:username, :taskid);", data)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
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
	ctx.HTML(http.StatusOK, "newtask.html", gin.H{ "Title"   : "NEW TASK",
																					       "Message" : message })
}
