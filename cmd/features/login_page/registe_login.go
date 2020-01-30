package login_page

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

type Account struct {
	c        *gin.Context
	username string
	password string
	uid      string
	time     string
}

func Start(c *gin.Context) {
	var a Account

	a.username = c.PostForm("username")
	//TODO:考虑用户安全，加密传输和存储用户密码
	a.password = c.PostForm("password")
	a.c = c
	//uid为用户注册的时间戳
	if a.IsPasswdExceedSix() {
		if a.IsRegiste("user") {
			a.Login()
		} else {
			a.Registe()
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

//判断密码是否大于六位
func (a Account) IsPasswdExceedSix() bool {
	return len(a.password) >= 6
}

//注册
func (a Account) Registe() {
	basic.G_UserID = basic.GetAUid()
	err := database.G_DB.Table.Insert("user", []string{"username", "password", "uid"}, []string{a.username, a.password, basic.G_UserID})
	basic.CheckError(err, "注册失败！")
	if err != nil {
		a.c.JSON(500, gin.H{
			"error": "注册失败！",
		})
	} else {
		a.c.JSON(http.StatusOK, gin.H{
			"message": "注册成功！",
		})
	}
}

//登录
func (a Account) Login() {
	var err error
	if a.HaveCookie("userID") {
		basic.G_UserID = basic.GetUid(a.c)
		a.c.JSON(http.StatusOK, gin.H{
			"message": "欢迎回来！" + a.username,
		})
	} else {
		if a.PasswdIsOk() {
			basic.G_UserID, err = database.UserName2Uid(a.username)
			basic.CheckError(err, "username查询uid失败!")
			basic.G_NickName, err = database.Uid2NickName(basic.G_UserID)
			basic.CheckError(err, "uid查询nickname失败!")
			a.c.SetCookie("userID", basic.G_UserID, 1000, "/", "127.0.0.1", false, true)
			a.c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "登录成功！",
			})
		} else {
			a.c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "密码错误",
					"code":    100003,
					"name":    "ERR_BAD_PASSWORD"},
			})
		}
	}
}

//判断密码是否正确
func (a Account) PasswdIsOk() bool {
	flag, err := database.G_DB.Table.Find("user", "password", "username", a.username)
	basic.CheckError(err, "判断密码是否正确失败！")
	return flag[0]["password"] != a.password
}

//判断是否已经注册
func (a Account) IsRegiste(tableName string) bool {
	flag, err := database.G_DB.Table.Find(tableName, "username", "username", a.username)
	basic.CheckError(err, "查询是否注册失败！")
	return flag != nil
}

//判断是否已有cookie来判断用户在线情况
func (a Account) HaveCookie(key string) bool {
	_, err := a.c.Cookie(key)
	basic.CheckError(err, "没有cookie！")
	return err == nil
}
