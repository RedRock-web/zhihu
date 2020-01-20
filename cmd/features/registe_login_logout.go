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

	if IsPasswdExceedSix(password) {
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

func IsPasswdExceedSix(passwd string) bool {
	return len(passwd) >= 6
}

func Registe(db *sql.DB, c *gin.Context, username string, password string) {
	if database.InsertField(db, "user_information", username, password) == nil{
		c.JSON(http.StatusOK, gin.H{
			"message": "注册成功！",
		})
	}
}

func Login(db *sql.DB, c *gin.Context, username string, password string) {
	if basic.HaveCookie(c, "userID") {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎回来！" + username,
		})
	} else {
		if PasswdIsOk(db, c, username, password) {
			//TODO:把cookie的value换成用户id提高安全，用户id使用username随机生成的字符串，加密算法还未定
			c.SetCookie("userID", username, 1000, "/", "127.0.0.1", false, true)
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
	user, _ := database.DatabaseSearch(db, "user_information", "username", username)
	return user.Password == password
}

func IsRegiste(db *sql.DB, tableName string, username string) bool {
	user, _ := database.DatabaseSearch(db, tableName, "username", username)
	return user.Id != ""
}

//fixme:chrome模拟注销成功，但是postman失败
func Logout(db *sql.DB, c *gin.Context) {
	cookie, _ := c.Cookie("userID")
	c.SetCookie("userID", cookie, -1, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusFound, gin.H{
		"message": "账号已注销！",
	})
}
