package login_page

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

func RegisteOrLogin(db *sql.DB, c *gin.Context, tableName string) string {
	username := c.PostForm("username")
	password := c.PostForm("password")
	uid := strconv.FormatInt(time.Now().Unix(), 10)

	if IsPasswdExceedSix(password) {
		if IsRegiste(db, tableName, username) {
			Login(db, c, username, password, uid)
		} else {
			Registe(db, c, username, password, uid)
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

func Registe(db *sql.DB, c *gin.Context, username string, password string, uid string) {
	if database.InsertField(db, "user", username, password, uid) == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "注册成功！",
		})
	} else {
		fmt.Println(errors.New("注册失败！"))
	}
}

func Login(db *sql.DB, c *gin.Context, username string, password string, uid string) {
	if basic.HaveCookie(c, "userID") {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎回来！" + username,
		})
	} else {
		if PasswdIsOk(db, c, username, password) {
			c.SetCookie("userID", uid, 1000, "/", "127.0.0.1", false, true)
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
	user, _ := database.DatabaseSearch(db, "user", "username", username)
	return user.Password == password
}

func IsRegiste(db *sql.DB, tableName string, username string) bool {
	user, _ := database.DatabaseSearch(db, tableName, "username", username)
	return user.Id != ""
}

func Logout(db *sql.DB, c *gin.Context) {
	cookie, _ := c.Cookie("userID")
	c.SetCookie("userID", cookie, -1, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusFound, gin.H{
		"message": "账号已注销！",
	})
}
