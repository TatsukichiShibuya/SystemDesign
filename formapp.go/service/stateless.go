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
    NewName string `form:"newname"`
    NewDate string `form:"newdate"`
    NewMessage string `form:"newmessage"`
    Name string `form:"name"`
    Date string `form:"date"`
    Message string `form:"message"`
}

func initData(data *FormData) {
    // initialize data with default values
    if data.Name == "" { data.Name = "" }
    if data.Date == "" { data.Date = "2000-01-01" }
    if data.Message == "" { data.Message = ""}
}

func saveData(data *FormData) {
    // if new values exist, replace old data
    if data.NewName != "" { data.Name = data.NewName }
    if data.NewDate != "" { data.Date = data.NewDate }
    if data.NewMessage != "" { data.Message = data.NewMessage}
}

func rootHandler(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "start.html", nil)
}

func nameHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)
    initData(&data)

		// save user's input
    saveData(&data)

    // show web page
    ctx.HTML(http.StatusOK, "name_form.html", &data)
}

func dateHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)
    initData(&data)

    // save user's input
    saveData(&data)

    // show web page
    ctx.HTML(http.StatusOK, "date_form.html", &data)
}

func messageHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)
    initData(&data)

    // save user's input
    saveData(&data)

    // show web page
    ctx.HTML(http.StatusOK, "message_form.html", &data)
}

func checkHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)
    initData(&data)

    // save user's input
    saveData(&data)

    // show web page
    ctx.HTML(http.StatusOK, "check_form.html", &data)
}
