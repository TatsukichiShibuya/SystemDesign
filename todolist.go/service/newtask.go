package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetNewtask(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	Newtask(ctx, "")
}

func PostNewtask(ctx *gin.Context) {
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

	// get inputs
	title, _ := ctx.GetPostForm("title")
	deadline_string, _ := ctx.GetPostForm("deadline")
	deadline := parseDeadline(deadline_string)
	share, _ := ctx.GetPostForm("share")
	shareusers := parseUsers(share)
	detail, _ := ctx.GetPostForm("detail")

	// check if title is given
	if title == "" {
		Newtask(ctx, "タスク名を記入してください")
		return
	}

	// check if all of shareusers exists
	if !allUsersExist(shareusers, db){
		Newtask(ctx, "「共有するユーザー」に存在しないユーザー名が指定されています")
		return
	}

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

	// get new task's id
	lastid, err := res.LastInsertId()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// set owner
	data = map[string]interface{}{ "userid": userid, "taskid": lastid }
	_, err = db.NamedExec("INSERT INTO owners (userid, taskid) VALUES (:userid, :taskid);", data)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// also share user
	var suser database.User
	for _, shareuser := range shareusers {
		// get shareuser's id
		err = db.Get(&suser, "SELECT * FROM users WHERE username=?", shareuser)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// add to owners
		data = map[string]interface{}{ "userid": suser.ID, "taskid": lastid }
		_, err = db.NamedExec("INSERT INTO owners (userid, taskid) VALUES (:userid, :taskid);", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	ctx.Redirect(http.StatusSeeOther, "/list")
}

func Newtask(ctx *gin.Context, message string) {
	ctx.HTML(http.StatusOK, "newtask.html", gin.H{ "Title"  : "NEW TASK",
																					       "Message": message })
}
