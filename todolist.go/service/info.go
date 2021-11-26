package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetInfo(ctx *gin.Context) {
	if sessionCheck(ctx) {
		Info(ctx, "")
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostInfo(ctx *gin.Context) {
	if sessionCheck(ctx) {
		// Get DB connection
		db, err := database.GetConnection()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// get inputs
		submit, _ := ctx.GetPostForm("submit")
		oldname, _ := ctx.GetPostForm("oldname")
		newname, _ := ctx.GetPostForm("newname")
		oldpass, _ := ctx.GetPostForm("oldpass")
		newpass, _ := ctx.GetPostForm("newpass")
		message := ""

		// check if oldname is correct
		session := sessions.Default(ctx)
		username := session.Get("username")
		if username != oldname {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		// get user info
		var user database.User
		err = db.Get(&user, "SELECT * FROM users WHERE username=?", username)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// check passward
		if user.Passward == hash(oldpass) {
			if submit == "delete" {
				// delete user
				data := map[string]interface{}{ "username": username }
				_, err = db.NamedExec("DELETE FROM users WHERE username=:username", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}

				// delete owner
				_, err = db.NamedExec("DELETE FROM owners WHERE username=:username", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}

				// logout
			  session.Delete("username")
			  session.Save()
				ctx.Redirect(http.StatusSeeOther, "/login")
				return
			} else if submit == "update" {
				// check if new info is not blank
				if newname == "" && newpass == "" {
					message = "更新情報を入力してください"
				} else {
					// update passward
					if newpass != "" {
						data := map[string]interface{}{ "username": username, "newpass": hash(newpass) }
						_, _ = db.NamedExec("UPDATE users SET passward=:newpass WHERE username=:username", data)
						message = "パスワードを変更しました"
					}

					// update username
					if newname != "" {
						// check if username is not used
						err = db.Get(&user, "SELECT * FROM users WHERE username=?", newname)
						if err != nil {
							// change username
							data := map[string]interface{}{ "oldname": oldname, "newname": newname }
							_, err = db.NamedExec("UPDATE users SET username=:newname WHERE username=:oldname", data)
							if err != nil {
								ctx.String(http.StatusInternalServerError, err.Error())
								return
							}
							session.Set("username", newname)
							session.Save()

							// update owners
							_, err = db.NamedExec("UPDATE owners SET username=:newname WHERE username=:oldname", data)
							if err != nil {
								ctx.String(http.StatusInternalServerError, err.Error())
								return
							}
							if message != "" {
								message = "ユーザー名とパスワードを変更しました"
							} else {
								message = "ユーザー名を変更しました"
							}
						} else {
							message = "指定されたユーザー名はすでに使用されています"
						}
					}
				}
			}
		} else {
			message = "パスワードが間違っています"
		}

		Info(ctx, message)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func Info(ctx *gin.Context, message string) {
	session := sessions.Default(ctx)
	username := session.Get("username")
	ctx.HTML(http.StatusOK, "info.html", gin.H{ "Title"    : "USER INFO",
																						  "Username" : username,
																					    "Message"  : message })
}
