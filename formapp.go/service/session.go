package main

import (
	"fmt"
	"net/http"
  "github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FormData struct {
		Name string `form:"name"`
		Date string `form:"date"`
		Message string `form:"message"`
}
var dataMap map[string]*FormData
const DEFAULT_NAME = "NONAME"
const DEFAULT_DATE = "2000-01-01"
const DEFAULT_MESSAGE = ""

// config
const port = 8000
const COOKIE_KEY = "id"

func main() {
		// initialize HashMap
		dataMap = make(map[string]*FormData)

    // initialize Gin engine
    engine := gin.Default()
		engine.LoadHTMLGlob("templates/session/*.html")

    // routing
    engine.GET("/", rootHandler)
		engine.POST("/", registerHandler)

    engine.POST("/name-form", nameHandler)
    engine.POST("/date-form", dateHandler)
    engine.POST("/message-form", messageHandler)
    engine.POST("/check-form", checkHandler)

    // start server
    engine.Run(fmt.Sprintf(":%d", port))
}

func getCookie(ctx *gin.Context) string {
	cookie, err := ctx.Request.Cookie(COOKIE_KEY)
	if err != nil {
		// 新しいセッション
		id, _ := uuid.NewRandom()
		ctx.SetCookie(COOKIE_KEY, id.String(), 600, "/", "localhost", false, true)
		dataMap[id.String()] = &FormData{ Name: DEFAULT_NAME,
															 			  Date: DEFAULT_DATE,
															 		    Message: DEFAULT_MESSAGE }
		return id.String()
	} else if _, exist := dataMap[cookie.Value]; !exist {
		// クッキーが存在するがデータがない
		//（サーバーが落ちてクライアント側のクッキーがまだ有効な状態）
		dataMap[cookie.Value] = &FormData{ Name: DEFAULT_NAME,
															 			  Date: DEFAULT_DATE,
															 		    Message: DEFAULT_MESSAGE }
	}
	return cookie.Value
}

func saveData(cookie string, data FormData) {
		if data.Name != "" {
			dataMap[cookie].Name = data.Name
		}
		if data.Date != "" {
			dataMap[cookie].Date = data.Date
		}
		if data.Message != "" {
			dataMap[cookie].Message = data.Message
		}
}


func rootHandler(ctx *gin.Context) {
		_ = getCookie(ctx)
		ctx.HTML(http.StatusOK, "start.html", nil)
}

func nameHandler(ctx *gin.Context) {
		// get data
    var data FormData
    _ = ctx.Bind(&data)

		// save data
		cookie := getCookie(ctx)
		saveData(cookie, data)

		// show page
		ctx.HTML(http.StatusOK, "name_form.html", dataMap[cookie])
}

func dateHandler(ctx *gin.Context) {
		// get data
		var data FormData
		_ = ctx.Bind(&data)

		// save data
		cookie := getCookie(ctx)
		saveData(cookie, data)

		// show page
		ctx.HTML(http.StatusOK, "date_form.html", dataMap[cookie])
}

func messageHandler(ctx *gin.Context) {
		// get data
		var data FormData
		_ = ctx.Bind(&data)

		// save data
		cookie := getCookie(ctx)
		saveData(cookie, data)

		// show page
		ctx.HTML(http.StatusOK, "message_form.html", dataMap[cookie])
}

func checkHandler(ctx *gin.Context) {
		// get data
		var data FormData
		_ = ctx.Bind(&data)

		// save data
		cookie := getCookie(ctx)
		saveData(cookie, data)

		// show page
    ctx.HTML(http.StatusOK, "check_form.html", dataMap[cookie])
}

func registerHandler(ctx *gin.Context) {
		// 登録処理をする.ここではデータ削除
		cookie := getCookie(ctx)
		delete(dataMap, cookie)
		ctx.HTML(http.StatusOK, "start.html", nil)
}
