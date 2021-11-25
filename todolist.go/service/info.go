package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

// Info renders info.html
func GetInfo(ctx *gin.Context) {
	if sessionCheck(ctx) {
		Info(ctx, "")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostInfo(ctx *gin.Context) {
	if sessionCheck(ctx) {
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		submit, _ := ctx.GetPostForm("submit")
		oldname, _ := ctx.GetPostForm("oldname")
		newname, _ := ctx.GetPostForm("newname")
		oldpass, _ := ctx.GetPostForm("oldpass")
		newpass, _ := ctx.GetPostForm("newpass")
		message := ""

		if submit == "delete" {
			var user database.User
			err := db.Get(&user, "SELECT * FROM users WHERE username=?", oldname)
			if err != nil {
				fmt.Println("err")
				return
			}
			// check passward
			if user.Passward == hash(oldpass) {
				// delete user
				data := map[string]interface{}{ "username": oldname }
				_, _ = db.NamedExec("DELETE FROM users WHERE username=:username", data)
				// delete owner
				_, _ = db.NamedExec("DELETE FROM owners WHERE username=:username", data)
				// logout
				session := sessions.Default(ctx)
			  session.Delete("username")
			  session.Save()
				ctx.Redirect(http.StatusSeeOther, "/login")
				return
			} else {
				message = "パスワードが間違っています"
			}
		} else if submit == "update" {
			var user database.User
			err := db.Get(&user, "SELECT * FROM users WHERE username=?", oldname)
			if err != nil {
				fmt.Println("err")
				return
			}
			// check passward
			if user.Passward == hash(oldpass) {
				if newname == "" && newpass == "" {
					message = "更新情報を入力してください"
				} else {
					if newname != "" {
						// check if username is not used
						err = db.Get(&user, "SELECT * FROM users WHERE username=?", newname)
						if err != nil {
							// change username
							data := map[string]interface{}{ "oldname": oldname, "newname": newname}
							_, _ = db.NamedExec("UPDATE users SET username=:newname WHERE username=:oldname", data)
							session := sessions.Default(ctx)
						  session.Set("username", newname)
							session.Save()
							// update owners
							_, _ = db.NamedExec("UPDATE owners SET username=:newname WHERE username=:oldname", data)
							message += "名前を変更しました"
						} else {
							message = "ユーザー名はすでに使用されています"
						}
					}
					if newpass != "" {
						// change username
						data := map[string]interface{}{ "oldname": oldname, "newpass": hash(newpass)}
						_, _ = db.NamedExec("UPDATE users SET passward=:newpass WHERE username=:oldname", data)
						if message != "" { message += ". " }
						message += "パスワードを変更しました"
					}
				}
			} else {
				message = "パスワードが間違っています"
			}
		} else {
			fmt.Println("err")
			return
		}

		if message == "" {
			fmt.Println("err")
			return
		} else {
			Info(ctx, message)
		}
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func Info(ctx *gin.Context, message string) {
	session := sessions.Default(ctx)
	username := session.Get("username")
	ctx.HTML(http.StatusOK, "info.html", gin.H{ "Title": "USER INFO",
																						  "Username": username,
																					    "Message": message })
}
