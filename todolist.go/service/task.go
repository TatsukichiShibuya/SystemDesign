package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetTask(ctx *gin.Context) {
	if sessionCheck(ctx) {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// parse ID given as a parameter
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		// check ownership
		session := sessions.Default(ctx)
		username := session.Get("username").(string)
		var owner database.Owner
		err = db.Get(&owner, "SELECT * FROM owners WHERE taskid=? AND username=?", id, username)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		// Get a task with given ID
		var task database.Task
		err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		Task(ctx, task, "")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostTask(ctx *gin.Context) {
	if sessionCheck(ctx) {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// parse ID given as a parameter
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		// check ownership
		session := sessions.Default(ctx)
		username := session.Get("username").(string)
		var owner database.Owner
		err = db.Get(&owner, "SELECT * FROM owners WHERE taskid=? AND username=?", id, username)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		submit, _ := ctx.GetPostForm("submit")
		message := ""

		if submit == "update" {
			// get parameters
			title, _ := ctx.GetPostForm("title")
			isdone, _ := ctx.GetPostForm("isdone")
			deadline_string, _ := ctx.GetPostForm("deadline")
			share, _ := ctx.GetPostForm("share")
			shareusers := parseUsers(share)
			detail, _ := ctx.GetPostForm("detail")

			if title == "" {
				message = "タイトルを入力してください"
			} else if !allExist(shareusers, db){
			 	message = "「共有しているユーザー」に存在しないユーザー名が指定されています"
		 	} else {
				// parse deadline
				deadline := parseDeadline(deadline_string)

				// update infomations with parameters
				data := map[string]interface{}{ "id": id, "title": title, "deadline": deadline, "detail": detail }
				if deadline.IsZero() {
					if isdone == "done" {
						_, err = db.NamedExec("UPDATE tasks SET title=:title, is_done=b'1', detail=:detail WHERE id=:id", data)
					} else if isdone == "undone" {
						_, err = db.NamedExec("UPDATE tasks SET title=:title, is_done=b'0', detail=:detail WHERE id=:id", data)
					}
				} else {
					if isdone == "done" {
						_, err = db.NamedExec("UPDATE tasks SET title=:title, is_done=b'1', deadline=:deadline, detail=:detail WHERE id=:id", data)
					} else if isdone == "undone" {
						_, err = db.NamedExec("UPDATE tasks SET title=:title, is_done=b'0', deadline=:deadline, detail=:detail WHERE id=:id", data)
					}
				}
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}

				// update ownership
				data = map[string]interface{}{ "taskid": id }
				_, err = db.NamedExec("DELETE FROM owners WHERE taskid=:taskid", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}
				for _, shareuser := range shareusers {
					data = map[string]interface{}{ "username": shareuser, "taskid": id }
					_, err = db.NamedExec("INSERT INTO owners (username, taskid) VALUES (:username, :taskid);", data)
					if err != nil {
						ctx.String(http.StatusInternalServerError, err.Error())
						return
					}
				}
				data = map[string]interface{}{ "username": username, "taskid": id }
				_, err = db.NamedExec("INSERT INTO owners (username, taskid) VALUES (:username, :taskid);", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}

				message += "タスクを更新しました"
			}
		}

		// Get a task with given ID
		var task database.Task
		err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id)
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		Task(ctx, task, message)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func Task(ctx *gin.Context, task database.Task, message string) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	session := sessions.Default(ctx)
	username := session.Get("username").(string)
	ftask, err := formatTask(task, ctx, db, username)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "task.html", gin.H{ "Title"   : "TASK",
																						  "Task"    : ftask,
																						  "Message" : message })
}