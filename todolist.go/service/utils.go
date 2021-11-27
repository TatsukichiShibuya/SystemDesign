package service

import (
	"fmt"
	"time"
	"strings"
	"unicode"
	"crypto/sha256"

	"github.com/gin-gonic/gin"
  "github.com/gin-contrib/sessions"
	"github.com/jmoiron/sqlx"

	database "todolist.go/db"
)

func sessionCheck(ctx *gin.Context) bool {
  session := sessions.Default(ctx)
  return session.Get("userid") != nil
}

func hash(passward string) string {
	h := sha256.Sum256([]byte(passward))
	return fmt.Sprintf("%x", h)
}

func isAcceptableString(str string) bool {
	res := true
	for _, c := range str {
		ok := false
		ok = ok || unicode.In(c, unicode.Hiragana)
		ok = ok || unicode.In(c, unicode.Katakana)
		ok = ok || unicode.In(c, unicode.Han)
		ok = ok || unicode.IsDigit(c)
		ok = ok || unicode.IsLetter(c)
		res = res && ok
	}
	return res
}

func isOwnerOfTask(userid uint64, taskid uint64, db *sqlx.DB) bool {
	var temp database.Owner
	err := db.Get(&temp, "SELECT * FROM owners WHERE userid=? AND taskid=?", userid, taskid)
	return err == nil
}

func parseUsers(users string) []string {
	var res []string
	if users == "" {
		return res
	} else {
		return strings.Split(users, ",")
	}
}

func parseDeadline(deadline string) time.Time {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	t, _ := time.ParseInLocation("2006-01-02T15:04", deadline, jst)
	return t
}

func allUsersExist(users []string, db *sqlx.DB) bool {
	var user database.User
	for _, username := range users {
		if isAcceptableString(username) {
			err := db.Get(&user, "SELECT * FROM users WHERE username=?", username)
			if err != nil {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func culcLimit(deadline time.Time) string {
	negative, days, hours, minutes, _ := parseTime(deadline)
	var res string
	if days >= 365 {
		res = fmt.Sprintf("%d年以上", int(days/365))
	} else if days >= 10 {
		res = fmt.Sprintf("%d日", days)
	} else {
		if days > 1 {
			res = fmt.Sprintf("%d日", days)
			if hours > 0 {
				res += fmt.Sprintf("%d時間", hours)
			}
		} else {
			if hours > 0 {
				res = fmt.Sprintf("%d時間", hours)
				res += fmt.Sprintf("%d分", minutes)
			} else {
				res = fmt.Sprintf("%d分", minutes)
			}
		}
	}
	if negative {
		res += "前"
	}
	return res
}

func culcImportance(deadline time.Time) int {
	negative, days, _, _, _ := parseTime(deadline)
	if negative {
		return 3
	} else if days == 0 {
		return 2
	} else {
		return 1
	}
}

func parseTime(t time.Time) (bool, int, int, int, int) {
	sub := int(t.Sub(time.Now()).Seconds())
	negative := (sub<0)
	if negative {
		sub *= -1
	}

	seconds := sub%60
	sub -= seconds

	sub /= 60
	minutes := sub%60
	sub -= minutes

	sub /= 60
	hours := sub%24
	sub -= hours

	sub /= 24
	days := sub

	return negative, days, hours, minutes, seconds
}
