package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetInfo(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	Info(ctx, "")
}

func PostInfo(ctx *gin.Context) {
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
	submit, _ := ctx.GetPostForm("submit")
	newname, _ := ctx.GetPostForm("newname")
	oldpass, _ := ctx.GetPostForm("oldpass")
	newpass, _ := ctx.GetPostForm("newpass")
	message := ""

	// get user info
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE id=?", userid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// check passward
	if user.Passward != hash(oldpass) {
		Info(ctx, "パスワードが間違っています")
		return
	}

	if submit == "delete" {
		// delete user
		data := map[string]interface{}{ "userid": userid }
		_, err = db.NamedExec("DELETE FROM users WHERE id=:userid", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// delete owner
		_, err = db.NamedExec("DELETE FROM owners WHERE userid=:userid", data)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// logout
		ctx.Redirect(http.StatusSeeOther, "/logout")
		return
	} else if submit == "update" {
		// no info to update
		if newname == "" && newpass == "" {
			Info(ctx, "更新情報を入力してください")
			return
		}

		// update passward if needed
		if newpass != "" {
			data := map[string]interface{}{ "userid": userid, "newpass": hash(newpass) }
			_, err = db.NamedExec("UPDATE users SET passward=:newpass WHERE id=:userid", data)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			message = "パスワードを変更しました"
		}

		// update username if needed
		if newname != "" {
			// check if username is acceptable
			if !isAcceptableString(newname) {
				Info(ctx, "ユーザー名には（アルファベット，数字，ひらがな，カタカナ，漢字）のみ使えます")
				return
			}

			// check if username is not used
			var temp database.User
			err = db.Get(&temp, "SELECT * FROM users WHERE username=?", newname)
			if err == nil {
				Info(ctx, "指定されたユーザー名はすでに使用されています")
				return
			}

			// change username
			data := map[string]interface{}{ "userid": userid, "newname": newname }
			_, err = db.NamedExec("UPDATE users SET username=:newname WHERE id=:userid", data)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			if message != "" {
				message = "ユーザー名とパスワードを変更しました"
			} else {
				message = "ユーザー名を変更しました"
			}
		}
	}
	Info(ctx, message)
}

func Info(ctx *gin.Context, message string) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	session := sessions.Default(ctx)
	userid := session.Get("userid").(uint64)
	var user database.User
	err = db.Get(&user, "SELECT * FROM users WHERE id=?", userid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "info.html", gin.H{ "Title"   : "USER INFO",
																							"Username": user.Username,
																							"Message" : message })
}
