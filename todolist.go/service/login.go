package service

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	database "todolist.go/db"
)

// Login renders login.html
func GetLogin(ctx *gin.Context) {
	Login(ctx, "")
}

func PostLogin(ctx *gin.Context) {
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	submit, _ := ctx.GetPostForm("submit")
	username, _ := ctx.GetPostForm("username")
	passward, _ := ctx.GetPostForm("passward")
	message := ""

	if submit == "login" {
		 // check username and passward are correct
		 if username == "" || passward == "" {
			 message = "記入が不足しています"
		 } else {
			 var user database.User
			 err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
			 if err != nil {
				 message = "ユーザーが存在しません"
			 } else {
				 // check passward
				 if user.Passward != hash(passward) {
					 message = "パスワードが間違っています"
				 }
			 }
		 }
	} else if submit == "register" {
		// check username and passward are given
		if username == "" || passward == "" {
			message = "記入が不足しています"
		} else {
			// check if username is not used
			var user database.User
			err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
			if err != nil {
				// register new user
				data := map[string]interface{}{ "username": username, "passward": hash(passward) }
				_, err = db.NamedExec("INSERT INTO users (username, passward) VALUES (:username, :passward)", data)
				if err != nil {
					fmt.Println("era-")
				}
			} else {
				message = "ユーザー名はすでに使用されています"
			}
		}
	}else {
		fmt.Println("err")
		return
	}

	if message == "" {
		// move to /home
		session := sessions.Default(ctx)
	  session.Set("username", username)
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		// remain because something is wrong
		Login(ctx, message)
	}
}

func Login(ctx *gin.Context, message string) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{"Title"   : "LOGIN",
																							"Message" : message})
}
