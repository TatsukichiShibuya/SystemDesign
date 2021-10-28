package main

import (
    //"os"
    "fmt"
    "net/http"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type FormData struct {
    Name string `form:"name"`
    Date string `form:"date"`
    Message string `form:"message"`
}
var dataMap map[string]*FormData
const DEFAULT_NAME = ""
const DEFAULT_DATE = "2000-01-01"
const DEFAULT_MESSAGE = ""

// config
const port = 8000
const COOKIE_KEY = "id"
const DATA_DIR = "data/"

func main() {
    // initialize HashMap with local data
    dataMap = make(map[string]*FormData)
    // search json file in DATA_DIR
    files, _ := ioutil.ReadDir(DATA_DIR)
    for _, file := range files {
        filename := file.Name()
        ext := filepath.Ext(filename)
        if ext == ".json" {
            // read json file
            raw, _ := ioutil.ReadFile(DATA_DIR+filename)
            var data FormData
            json.Unmarshal(raw, &data)
            // register {name: &data}
            name := filename[0:len(filename)-len(ext)]
            dataMap[name] = &data
        }
    }

    // initialize Gin engine
    engine := gin.Default()
    engine.LoadHTMLGlob("templates/persistent/*.html")

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

func getCookie(ctx *gin.Context) string {
    cookie, err := ctx.Request.Cookie(COOKIE_KEY)
    if err != nil {
        fmt.Println("UP")
        // new session
        id, _ := uuid.NewRandom()
        ctx.SetCookie(COOKIE_KEY, id.String(), 600, "/", "localhost", false, true)
        dataMap[id.String()] = &FormData{ Name: DEFAULT_NAME,
                                          Date: DEFAULT_DATE,
                                          Message: DEFAULT_MESSAGE }
        saveJson(id.String())
        return id.String()
    } else if _, exist := dataMap[cookie.Value]; !exist {
        fmt.Println("DOWN")
        // cookie exists but haven't registered
        dataMap[cookie.Value] = &FormData{ Name: DEFAULT_NAME,
                                           Date: DEFAULT_DATE,
                                           Message: DEFAULT_MESSAGE }
        saveJson(cookie.Value)
    }
    return cookie.Value
}

func saveData(cookie string, data FormData) {
    // if new values exist, replace old data
    if data.Name != "" { dataMap[cookie].Name = data.Name }
    if data.Date != "" { dataMap[cookie].Date = data.Date }
    if data.Message != "" { dataMap[cookie].Message = data.Message }
    saveJson(cookie)
}

func saveJson(name string) {
    j, _ := json.Marshal(dataMap[name])
    ioutil.WriteFile(DATA_DIR+name+".json", []byte(j), 0664)
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

    // show web page
    ctx.HTML(http.StatusOK, "name_form.html", dataMap[cookie])
}

func dateHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)

    // save data
    cookie := getCookie(ctx)
    saveData(cookie, data)

    // show web page
    ctx.HTML(http.StatusOK, "date_form.html", dataMap[cookie])
}

func messageHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)

    // save data
    cookie := getCookie(ctx)
    saveData(cookie, data)

    // show web page
    ctx.HTML(http.StatusOK, "message_form.html", dataMap[cookie])
}

func checkHandler(ctx *gin.Context) {
    // get data
    var data FormData
    _ = ctx.Bind(&data)

    // save data
    cookie := getCookie(ctx)
    saveData(cookie, data)

    // show web page
    ctx.HTML(http.StatusOK, "check_form.html", dataMap[cookie])
}
