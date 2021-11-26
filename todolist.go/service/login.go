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
	message := ""

	// check if username and passward are given
	if username == "" || passward == "" {
		message = "記入が不足しています"
	} else {
		if submit == "login" {
			 // check username and passward are correct
			var user database.User
			err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
			if err != nil {
				message = "指定されたユーザー名は存在しません"
			} else {
				// check passward
				if user.Passward != hash(passward) {
				 message = "パスワードが間違っています"
				}
			}
		} else if submit == "register" {
			// check if username is not used
			var user database.User
			err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
			if err != nil {
				// register new user
				data := map[string]interface{}{ "username": username, "passward": hash(passward) }
				_, err = db.NamedExec("INSERT INTO users (username, passward) VALUES (:username, :passward)", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}
			} else {
				message = "指定されたユーザー名はすでに使用されています"
			}
		} else {
			ctx.String(http.StatusBadRequest, "不正なリクエスト")
			return
		}
	}

	if message == "" {
		// move to "/"
		session := sessions.Default(ctx)
	  session.Set("username", username)
		session.Options(sessions.Options{ MaxAge: 60000, })
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		// remain because something is wrong
		Login(ctx, message)
	}
}

func Login(ctx *gin.Context, message string) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{ "Title"   : "LOGIN",
																							 "Message" : message })
}
