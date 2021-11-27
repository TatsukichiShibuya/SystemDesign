package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")
	store := cookie.NewStore([]byte("secret"))
  engine.Use(sessions.Sessions("mysession", store))

	// routing
	engine.Static("/assets", "./assets")

	engine.GET("/", service.GetHome)
	engine.POST("/", service.GetHome)

	engine.GET("/login", service.GetLogin)
	engine.POST("/login", service.PostLogin)

	engine.GET("/info", service.GetInfo)
	engine.POST("/info", service.PostInfo)

	engine.GET("/list", service.GetList)
	engine.POST("/list", service.PostList)

	engine.GET("/task/:id", service.GetTask)
	engine.POST("/task/:id", service.PostTask)

	engine.GET("/newtask", service.GetNewtask)
	engine.POST("/newtask", service.PostNewtask)

	engine.GET("/logout", service.GetLogout)
	engine.POST("/logout", service.GetLogout)

	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
