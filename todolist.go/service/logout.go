package service

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

func GetLogout(ctx *gin.Context) {
	if sessionCheck(ctx) {
	  session := sessions.Default(ctx)
	  session.Delete("username")
	  session.Save()
	}
	ctx.Redirect(http.StatusSeeOther, "/login")
}
