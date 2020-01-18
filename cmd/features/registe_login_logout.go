package features

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

func RegisteOrLogin(db *sql.DB, c *gin.Context, tableName string) string {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if IsPasswdMoreSixDigit(password) {
		if IsRegiste(db, tableName, username) {
			Login(db, c, username, password)
		} else {
			Registe(db, c, username, password)
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "密码长度不足",
				"code":    100004,
				"name":    "ERR_BAD_PASSWORD_FORMAT"},
		})
	}

	return username
}

func IsPasswdMoreSixDigit(passwd string) bool {
	return len(passwd) >= 6
}

func Registe(db *sql.DB, c *gin.Context, username string, password string) {
	database.InsertField(db, "user_information", username, password)
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
		if PasswdIsOk(db, c, username, password) {
			c.SetCookie(username, username, 100, "/", "127.0.0.1", false, true)
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "登录成功！",
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "密码错误",
					"code":    100003,
					"name":    "ERR_BAD_PASSWORD"},
			})
		}
	}
}

func PasswdIsOk(db *sql.DB, c *gin.Context, username string, password string) bool {
	user := database.DatabaseSearch(db, "user_information", "username", username)
	return user.Password == password
}

func IsRegiste(db *sql.DB, tableName string, username string) bool {
	user := database.DatabaseSearch(db, tableName, "username", username)

	return user.Username != ""
}

//暂时出了点问题，无法删除cookie
func Logout(db *sql.DB, c *gin.Context, username string) {
	_, err := c.Cookie(username)
	//fmt.Println(username)
	if err == nil {
		c.SetCookie(username, username, -1, "/", "127.0.0.1", false, true)
		c.JSON(http.StatusFound, gin.H{
			"message": "账号已登出！",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "您并未登录，无需登出！",
		})
	}
}
