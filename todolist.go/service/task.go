package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetTask(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(ctx)
	userid := session.Get("userid").(uint64)

	// parse ID given as a parameter
	taskid_int, err := strconv.Atoi(ctx.Param("id"))
	taskid := uint64(taskid_int)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// check ownership
	var owner database.Owner
	err = db.Get(&owner, "SELECT * FROM owners WHERE taskid=? AND userid=?", taskid, userid)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// get a task whose id is given
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", taskid)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	Task(ctx, task, "")
}

func PostTask(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(ctx)
	userid := session.Get("userid").(uint64)

	// parse ID given as a parameter
	taskid_int, err := strconv.Atoi(ctx.Param("id"))
	taskid := uint64(taskid_int)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// check ownership
	var owner database.Owner
	err = db.Get(&owner, "SELECT * FROM owners WHERE taskid=? AND userid=?", taskid, userid)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// get a task whose id is given
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", taskid)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// get inputs
	submit, _ := ctx.GetPostForm("submit")
	title, _ := ctx.GetPostForm("title")
	isdone, _ := ctx.GetPostForm("isdone")
	deadline_string, _ := ctx.GetPostForm("deadline")
	deadline := parseDeadline(deadline_string)
	share, _ := ctx.GetPostForm("share")
	shareusers := parseUsers(share)
	detail, _ := ctx.GetPostForm("detail")

	// check if title is given
	if title == "" {
		Task(ctx, task, "タイトルを入力してください")
		return
	}

	// check if all of shareusers exists
	if !allUsersExist(shareusers, db){
		Task(ctx, task, "「共有するユーザー」に存在しないユーザー名が指定されています")
		return
	}

	if submit == "update" {
		// update infomations with parameters
		data := map[string]interface{}{ "id": taskid, "title": title, "deadline": deadline, "detail": detail }
		options := ""
		if isdone == "done" {
			options += " is_done=b'1',"
		} else if isdone == "undone" {
			options += " is_done=b'0',"
		}
		if !deadline.IsZero() {
			options += " deadline=:deadline,"
		}
		_, err = db.NamedExec("UPDATE tasks SET title=:title," + options + " detail=:detail WHERE id=:id", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// update ownership
		data = map[string]interface{}{ "taskid": taskid }
		_, err = db.NamedExec("DELETE FROM owners WHERE taskid=:taskid", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		data = map[string]interface{}{ "userid": userid, "taskid": taskid }
		_, err = db.NamedExec("INSERT INTO owners (userid, taskid) VALUES (:userid, :taskid);", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		for _, shareuser := range shareusers {
			// get shareuser's id
			var suser database.User
			err = db.Get(&suser, "SELECT * FROM users WHERE username=?", shareuser)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}

			// add to owners
			data = map[string]interface{}{ "userid": suser.ID, "taskid": taskid }
			_, err = db.NamedExec("INSERT INTO owners (userid, taskid) VALUES (:userid, :taskid);", data)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	// Get updated task with given ID
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", taskid)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	Task(ctx, task, "タスクを更新しました")
}

func Task(ctx *gin.Context, task database.Task, message string) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(ctx)
	userid := session.Get("userid").(uint64)

	// format task
	ftask, err := formatTask(task, db, userid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "task.html", gin.H{ "Title"  : "TASK",
																						  "Task"   : ftask,
																						  "Message": message })
}
