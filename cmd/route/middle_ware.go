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

//点击登录注册时,用cookie判断是否已经登录,如果已经登录,则重定向到主页,无法再次登录注册
func Unauthorized2LoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if features.IsLogin(c, "userID") {
			basic.RediRect(c, "/")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

//点击登录后才能点的url,用cookie判断是否已经登录,否则重定向到登录注册页
func Authorized2Some() gin.HandlerFunc {
	return func(c *gin.Context) {
		if features.IsLogin(c, "userID") {
			features.G_user.Info.Uid, _ = c.Cookie("userID")
			c.Next()
		} else {
			basic.RediRect(c, "/sign_in")
			c.Abort()
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

//判断是否回答了问题
func JudgeIfReply() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//判断是否关注了问题
func JudgeIfFollowQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := features.NewQuestion()
		q.Id = c.Param("questionId")
		if q.IsFollow() && basic.MethodIsOk(c, "DELETE") {
			if q.CancelFollow() == nil {
				c.JSON(http.StatusOK, gin.H{
					"isFollow": "no",
				})
			}
		} else if !q.IsFollow() && basic.MethodIsOk(c, "POST") {
			if q.Follow() == nil {
				c.JSON(http.StatusOK, gin.H{
					"isFollow": "yes",
				})
			}
		}
	}
}

//判断是否评论是该回答的评论
func JudgeIfCommentInAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//判断是否评论是该问题的评论
func JudgeIfCommentInQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
