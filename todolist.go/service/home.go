package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHome(ctx *gin.Context) {
	if sessionCheck(ctx) {
		Home(ctx)
	} else {
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func PostHome(ctx  *gin.Context) {
	GetHome(ctx)
}

func Home(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home.html", gin.H{ "Title" : "HOME" })
}
