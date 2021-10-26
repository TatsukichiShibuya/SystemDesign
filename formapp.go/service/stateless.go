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
    engine.LoadHTMLGlob("templates/stateless/*.html")

    // routing
    engine.GET("/", rootHandler)
    engine.POST("/", rootHandler)

    engine.POST("/name-form", nameHandler)
    engine.POST("/date-form", dateHandler)
    engine.POST("/message-form", messageHandler)
    engine.POST("/check-form", checkHandler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

type FormData struct {
    Name string `form:"name"`
    Date string `form:"date"`
    Message string `form:"message"`
    DefaultName string `form:"name"`
    DefaultDate string `form:"date"`
    DefaultMessage string `form:"message"`
}

func rootHandler(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "start.html", nil)
}

func nameHandler(ctx *gin.Context) {
    var data FormData
    _ = ctx.Bind(&data)

    ctx.HTML(http.StatusOK, "name_form.html", &data)
}

func dateHandler(ctx *gin.Context) {
    var data FormData
    _ = ctx.Bind(&data)

    ctx.HTML(http.StatusOK, "date_form.html", &data)
}

func messageHandler(ctx *gin.Context) {
    var data FormData
    _ = ctx.Bind(&data)

    ctx.HTML(http.StatusOK, "message_form.html", &data)
}

func checkHandler(ctx *gin.Context) {
    var data FormData
    _ = ctx.Bind(&data)

    ctx.HTML(http.StatusOK, "check_form.html", &data)
}
