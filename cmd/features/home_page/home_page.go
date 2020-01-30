package home_page

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
)

//注销
func Logout(c *gin.Context) {
	cookie, err := c.Cookie("userID")
	basic.CheckError(err, "注销失败！")
	c.SetCookie("userID", cookie, -1, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusFound, gin.H{
		"message": "账号已注销！",
	})
}
