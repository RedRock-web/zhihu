package registeOrLogin

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

func RegisteOrLogin(db *sql.DB, c *gin.Context, tableName string) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if IsRegiste(db, "user", username) {
		Login(db, c, username, password)
	} else {
		Registe(db, c, username, password)
	}
}

func IsPasswdMoreSixDigit(passwd string) bool {
	return len(passwd) >= 6
}

func Registe(db *sql.DB, c *gin.Context, username string, password string) {
	database.InsertField(db, "user", username, password)
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功！",
	})
}



func Login(db *sql.DB, c *gin.Context, username string, password string) {
	if basic.HaveCookie(c, username) {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎回来！",
		})
	} else {
		if IsPasswdMoreSixDigit(password) {
			if PasswdIsOk(db, c, username, password) {
				c.SetCookie(username, username, 100, "/", "127.0.0.1", false, true)
				c.JSON(http.StatusOK, gin.H{
					"status":  http.StatusOK,
					"message": "登录成功！",
				})
			} else {
				c.JSON(401, gin.H{
					"error": gin.H{"message": "密码错误",
						"code": 100003, "name": "ERR_BAD_PASSWORD"},
				})
			}
		} else {
			c.JSON(401, gin.H{
				"error": gin.H{"message": "密码长度不足",
					"code": 100004, "name": "ERR_BAD_PASSWORD_FORMAT"},
			})
		}
	}

}

func PasswdIsOk(db *sql.DB, c *gin.Context,username string, password string) bool {
	_, passwd := database.DatabaseSearch(db, "user", username)
	return passwd == password
}



func IsRegiste(db *sql.DB, tableName string, username string) bool {
	var (
		id     string
		name   string
		passwd string
		judge  bool
	)

	selectOder := "select * from " + tableName + " where username= \"" + username + "\""
	stmt, err := db.Query(selectOder)
	basic.CheckError(err)
	for stmt.Next() {
		stmt.Scan(&id, &name, &passwd)
	}
	if name != "" {
		judge = true
	} else {
		judge = false
	}

	return judge
}
