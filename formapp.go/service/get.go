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
    engine.LoadHTMLGlob("templates/get/*.html")

    // routing
    engine.GET("/", rootHandler)

    engine.GET("/sample", sampleHandler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

type Data struct {
    Amount string `form:"amount"`
}

func rootHandler(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "start.html", nil)
}

func sampleHandler(ctx *gin.Context) {
    // get data
    var data Data
    _ = ctx.Bind(&data)

    // show web page
    ctx.HTML(http.StatusOK, "sample.html", &data)
}
