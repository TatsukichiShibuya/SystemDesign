package service

import (
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"

	database "todolist.go/db"
)

func GetList(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(ctx)
	userid := session.Get("userid").(uint64)

	// get inputs
	title, _ := ctx.GetQuery("title")
	isdone, _ := ctx.GetQuery("isdone")
	deadline, _ := ctx.GetQuery("deadline")
	// other options (deadline and sort) are processed
	// using formatTasksWithOption in List

	query := "SELECT * FROM tasks WHERE id in (SELECT taskid FROM owners WHERE userid=?) AND title LIKE ?"
	if isdone == "done" {
		query += " AND is_done=b'1'"
	} else if isdone == "undone" {
		query += " AND is_done=b'0'"
	}
	if deadline == "yes" {
		query += " AND deadline > created_at"
	} else if deadline == "no" {
		query += " AND deadline <= created_at"
	}

	// Get tasks from DB
	var tasks []database.Task
	args := []interface{}{ userid, "%"+title+"%" }
	err = db.Select(&tasks, query, args...)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	List(ctx, tasks)
}

type CheckBox struct {
	CheckedIDs []string `form:"check[]"`
}

func PostList(ctx *gin.Context) {
	// check if the user is logged in
	if !sessionCheck(ctx) {
		ctx.Redirect(http.StatusSeeOther, "/login")
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(ctx)
	userid := session.Get("userid").(uint64)

	// get inputs
	submit, _ := ctx.GetPostForm("submit")
	var checkbox CheckBox
 	ctx.Bind(&checkbox)

	if submit == "complete" {
		// make "isdone" done
		for i:=0; i<len(checkbox.CheckedIDs); i++ {
			checkedid_int, err := strconv.Atoi(checkbox.CheckedIDs[i])
			checkedid := uint64(checkedid_int)
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if isOwnerOfTask(userid, checkedid, db) {
				data := map[string]interface{}{ "id": checkedid }
				_, err = db.NamedExec("UPDATE tasks SET is_done=b'1' WHERE id=:id", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
	} else if submit == "delete" {
		// delete tasks
		for i:=0; i<len(checkbox.CheckedIDs); i++ {
			checkedid_int, err := strconv.Atoi(checkbox.CheckedIDs[i])
			checkedid := uint64(checkedid_int)
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if isOwnerOfTask(userid, checkedid, db) {
				data := map[string]interface{}{ "taskid": checkedid }
				_, err = db.NamedExec("DELETE FROM owners WHERE taskid=:taskid", data)
				if err != nil {
					ctx.String(http.StatusInternalServerError, err.Error())
					return
				}
			}
		}
	}

	GetList(ctx)
}

func List(ctx *gin.Context, tasks []database.Task) {
	// format tasks
	ftasks, err := formatTasksWithOption(tasks, ctx)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "task_list.html", gin.H{ "Title": "TASK LIST",
																									 "Tasks": ftasks })
}
