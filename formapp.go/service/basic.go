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
    engine.LoadHTMLGlob("templates/basic/*.html")

    authorized := engine.Group("/", gin.BasicAuth(gin.Accounts{
        "user1": "pass",
        "user2": "word",
    }))

    // routing
    engine.GET("/", rootHandler)
    authorized.GET("/secret", secretHandler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

func rootHandler(ctx *gin.Context) {
    // show web page
    ctx.HTML(http.StatusOK, "start.html", nil)
}

func secretHandler(ctx *gin.Context) {
    // show web page
    ctx.HTML(http.StatusOK, "secret.html", nil)
}
