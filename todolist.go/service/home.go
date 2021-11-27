package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHome(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	Home(ctx)
}

func Home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", gin.H{ "Title": "HOME" })
}
