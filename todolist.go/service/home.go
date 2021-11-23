package service

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

// Home renders home.html
func GetHome(ctx *gin.Context) {
	if sessionCheck(ctx) {
		Home(ctx)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostHome(ctx *gin.Context) {
	if sessionCheck(ctx) {
		submit, _ := ctx.GetPostForm("submit")
		if submit == "logout" {
		  session := sessions.Default(ctx)
		  session.Delete("username")
		  session.Save()
			ctx.Redirect(http.StatusSeeOther, "/login")
			return
		} else {
			fmt.Println("err")
			return
		}
		Home(ctx)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	username := session.Get("username")
	ctx.HTML(http.StatusOK, "home.html", gin.H{"Title": "HOME",
																						 "Username": username})
}
