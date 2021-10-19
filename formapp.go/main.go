package main

import (
	"fmt"
	"net/http"
  "github.com/gin-gonic/gin"
)

// config
const port = 8000

func main() {
    // initialize Gin engine
    engine := gin.Default()
		engine.LoadHTMLGlob("templates/*.html")

    // routing
    engine.GET("/", rootHandler)
    engine.GET("/name-form", nameFormHandler)
    engine.POST("/register-name", registerNameHandler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

func rootHandler(ctx *gin.Context) {
    //ctx.String(http.StatusOK, "Hello world.")
		ctx.HTML(http.StatusOK, "hello.html", nil)
}

func nameFormHandler(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "name_form.html", nil)
}

type FormData struct {
		Name string `form:"name"`
}

func registerNameHandler(ctx *gin.Context) {
		var data FormData
		_ = ctx.Bind(&data)
		ctx.HTML(http.StatusOK, "result.html", &data)
}
