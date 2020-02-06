package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/features"
)


// 处理跨域请求,支持options访问的全局middleWare
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

//通过cookie判断是否进入登录注册页的middleWare
func LoginPageJudgeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if features.IsLogin(c, "userID") {
			basic.RediRect(c, "/")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

//使用cookie检测是否登录的middleWare
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("userID")
		if cookie == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登录！"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

//404 响应middleWare
func NoResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error_code": 50,
			"error_msg":  "page not exists",
		})
	}
}

//登录后,只要有新的请求,便刷新cookie的middleWare
func RefreshCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("userID")
		if err == nil {
			c.SetCookie("userID", cookie, 1000, "/", "127.0.0.1", false, true)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登录！"})
			c.Abort()
		}
	}
}
