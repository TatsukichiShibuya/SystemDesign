package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetList(ctx *gin.Context) {
	if sessionCheck(ctx) {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Get username
		session := sessions.Default(ctx)
		username := session.Get("username")

		// Get Query Parameters
		title, _ := ctx.GetQuery("title")
		isdone, _ := ctx.GetQuery("isdone")

		query := "SELECT * FROM tasks WHERE id in (SELECT taskid FROM owners WHERE username=?)"
		query += fmt.Sprintf(" AND title LIKE '%%%s%%'", title)
		if isdone == "done" {
			query += " AND is_done=b'1'"
		} else if isdone == "undone" {
			query += " AND is_done=b'0'"
		}

		// Get tasks in DB
		var tasks []database.Task
		err = db.Select(&tasks, query, username)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		List(ctx, tasks)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

type CheckBox struct {
	CheckIDs []string `form:"check[]"`
}

func PostList(ctx *gin.Context) {
	if sessionCheck(ctx) {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		submit, _ := ctx.GetPostForm("submit")
		var checkbox CheckBox
	 	ctx.Bind(&checkbox)

		if submit == "complete" {
			for i:=0; i<len(checkbox.CheckIDs); i++ {
				id := checkbox.CheckIDs[i]
				data := map[string]interface{}{ "id": id }
				_, err = db.NamedExec("UPDATE tasks SET is_done=b'1' WHERE id=:id", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}
			}
		} else if submit == "delete" {
			for i:=0; i<len(checkbox.CheckIDs); i++ {
				id := checkbox.CheckIDs[i]
				data := map[string]interface{}{ "taskid": id }
				_, err = db.NamedExec("DELETE FROM owners WHERE taskid=:taskid", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}
			}
		}

		GetList(ctx)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func List(ctx *gin.Context, tasks []database.Task) {
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{ "Title" : "TASK LIST",
																									 "Tasks" : formatTasks(tasks) })
}
