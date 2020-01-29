package login_page

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

func Start(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	basic.G_UserID = strconv.FormatInt(time.Now().Unix(), 10)

	if IsPasswdExceedSix(password) {
		if IsRegiste("user", username) {
			Login(c, username, password)
		} else {
			Registe(c, username, password)
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "密码长度不足",
				"code":    100004,
				"name":    "ERR_BAD_PASSWORD_FORMAT"},
		})
	}
}

func IsPasswdExceedSix(passwd string) bool {
	return len(passwd) >= 6
}

func Registe(c *gin.Context, username string, password string) {
	if database.G_DB.Table.Insert("user", []string{"username", "password", "uid"}, []string{username, password, basic.G_UserID}) == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "注册成功！",
		})
	} else {
		fmt.Println(errors.New("注册失败！"))
	}
}

func Login(c *gin.Context, username string, password string) {
	var err error
	if HaveCookie(c, "userID") {
		c.JSON(http.StatusOK, gin.H{
			"message": "欢迎回来！" + username,
		})
	} else {
		if PasswdIsOk(username, password) {
			basic.G_UserID, err = database.UserName2Uid(username)
			basic.CheckError(err, "username查询uid失败!")
			basic.G_NickName, err = database.Uid2NickName(basic.G_UserID)
			basic.CheckError(err, "uid查询nickname失败!")
			c.SetCookie("userID", basic.G_UserID, 1000, "/", "127.0.0.1", false, true)
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

func PasswdIsOk(username string, password string) bool {
	flag, err := database.G_DB.Table.Find("user", "password", "username", username)
	basic.CheckError(err, "判断密码是否正确失败！")
	return flag[0]["password"] != password
}

func IsRegiste(tableName string, username string) bool {
	flag, err := database.G_DB.Table.Find(tableName, "username", "username", username)
	basic.CheckError(err, "查询是否注册失败！")
	return flag != nil
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("userID")
	basic.CheckError(err, "注销失败！")
	c.SetCookie("userID", cookie, -1, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusFound, gin.H{
		"message": "账号已注销！",
	})
}

func HaveCookie(c *gin.Context, key string) bool {
	_, err := c.Cookie(key)
	return err == nil
}
