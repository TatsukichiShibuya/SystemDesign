package service

import (
	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
)

func sessionCheck(ctx *gin.Context) bool {
  session := sessions.Default(ctx)
  return session.Get("username") != nil
}
