package features

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//Account 表示一个帐号
type Account struct {
	C        *gin.Context
	Username string //帐号
	Password string //密码
	Time     string //注册时间
}

//NewAccount 创建一个帐号对象
func NewAccount() *Account {
	return &Account{}
}

//登录注册接口
func RegisteOrLogin(c *gin.Context) {
	a := NewAccount()
	a.Username = c.PostForm("username")
	//TODO:考虑用户安全，加密传输和存储用户密码
	a.Password = c.PostForm("password")
	a.C = c
	if a.IsRegiste("user") {
		if a.Login() == nil {
			basic.Redirect(c, "/")
		}
	} else {
		if a.Registe() == nil {
			basic.Redirect(c, "/")
		}
	}
}

//判断密码是否大于六位
func (a Account) IsPasswdExceedSix() bool {
	return len(a.Password) >= 6
}

//注册
func (a Account) Registe() (err error) {
	G_user.Info.Uid = basic.GetAUid()
	G_user.Info.Nickname = "知乎用户"
	err = database.G_DB.Table.Insert("user", []string{"username", "password", "uid"}, []string{a.Username, a.Password, G_user.Info.Uid})
	basic.CheckError(err, "注册失败！")
	if err != nil {
		a.C.JSON(500, gin.H{
			"error": "注册失败！",
		})
	}
	return err
}

// 登录
func (a Account) Login() (err error) {
	//先判断密码是否大于六位
	if a.IsPasswdExceedSix() {
		//再判断密码是否正确
		if a.PasswdIsOk() {
			G_user.Info.Uid, err = database.UserName2Uid(a.Username)
			basic.CheckError(err, "username查询uid失败!")
			if err != nil {
				return err
			}
			G_user.Info.Nickname, err = database.Uid2NickName(G_user.Info.Uid)
			basic.CheckError(err, "uid查询nickname失败!")
			if err != nil {
				return err
			}
			a.C.SetCookie("userID", G_user.Info.Uid, 1000, "/", "127.0.0.1", false, true)
		} else {
			a.C.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "密码错误",
					"code":    100003,
					"name":    "ERR_BAD_PASSWORD"},
			})
		}
	} else {
		a.C.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "密码长度不足",
				"code":    100004,
				"name":    "ERR_BAD_PASSWORD_FORMAT"},
		})
	}
	return err
}

//判断密码是否正确
func (a Account) PasswdIsOk() bool {
	flag, err := database.G_DB.Table.Find("user", "password", "username", a.Username)
	basic.CheckError(err, "判断密码是否正确失败！")
	return string(flag[0]["password"].([]uint8)) == a.Password
}

//判断是否已经注册
func (a Account) IsRegiste(tableName string) bool {
	flag, err := database.G_DB.Table.Find(tableName, "username", "username", a.Username)
	basic.CheckError(err, "查询是否注册失败！")
	return flag != nil
}

//判断是否已经登录
func (a Account) IsLogin(key string) bool {
	k, _ := a.C.Cookie(key)
	return k != ""
}

//判断是否已经登录
func IsLogin(c *gin.Context, key string) bool {
	k, _ := c.Cookie(key)
	fmt.Println("###")
	fmt.Println(k != "")
	return k != ""
}

//注销帐号
func Logout(c *gin.Context) {
	result, err := c.Cookie("userID")
	basic.CheckError(err, "注销失败！")
	c.SetCookie("userID", result, -1, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusFound, gin.H{
		"message": "账号已注销！",
	})
}

//注销帐号
func (a Account) Logout() error {
	result, err := a.C.Cookie("userID")
	basic.CheckError(err, "注销失败！")
	a.C.SetCookie("userID", result, -1, "/", "127.0.0.1", false, true)
	a.C.JSON(http.StatusFound, gin.H{
		"message": "账号已注销！",
	})
	return err
}

//TODO:重置密码
func PasswdReset(c *gin.Context) {

}
