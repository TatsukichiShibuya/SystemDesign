package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetLogin(ctx *gin.Context) {
	Login(ctx, "")
}

func PostLogin(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// get inputs
	submit, _ := ctx.GetPostForm("submit")
	username, _ := ctx.GetPostForm("username")
	passward, _ := ctx.GetPostForm("passward")

	// check if username and passward are given
	if username == "" || passward == "" {
		Login(ctx, "記入が不足しています")
		return
	}

	// check if username is acceptable
	if !isAcceptableString(username) {
		Login(ctx, "ユーザー名には（アルファベット，数字，ひらがな，カタカナ，漢字）のみ使えます")
		return
	}

	if submit == "login" {
		 // check username and passward are correct
		var user database.User
		err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
		if err != nil {
			Login(ctx, "指定されたユーザー名は存在しません")
			return
		}
		if user.Passward != hash(passward) {
			Login(ctx, "パスワードが間違っています")
			return
		}
	} else if submit == "register" {
		// check if username is not used
		var temp database.User
		err = db.Get(&temp, "SELECT * FROM users WHERE username=?", username)
		if err == nil {
			Login(ctx, "指定されたユーザー名はすでに使用されています")
		}

		// register new user
		data := map[string]interface{}{ "username": username, "passward": hash(passward) }
		_, err = db.NamedExec("INSERT INTO users (username, passward) VALUES (:username, :passward)", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		Login(ctx, "不正なアクセスです")
		return
	}

	// move to "/"
	session := sessions.Default(ctx)
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	session.Set("userid", user.ID)
	session.Set("username", user.Username)
	session.Options(sessions.Options{ MaxAge: 60*60*24, }) // login remains for one day
	session.Save()
	ctx.Redirect(http.StatusSeeOther, "/")
}

func Login(ctx *gin.Context, message string) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{ "Title"  : "LOGIN",
																							 "Message": message })
}
