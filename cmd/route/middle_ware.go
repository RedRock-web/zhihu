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
		}
	}
}

//写回答中间件
func ReplyAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAnswer()
		a.Uid = features.G_user.Info.Uid
		a.Time = basic.GetTimeNow()
		a.Id = basic.GetAQuestionId()
		a.Content = c.PostForm("content")
		a.QuestionId = c.Param("questionId")
		if !a.HaveAnswer() && basic.MethodIsOk(c, "POST") && a.Post() == nil {
			//
		}
		c.Abort()
	}
}

//删回答中间件
func DeleteAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAnswer()
		a.Uid = features.G_user.Info.Uid
		a.QuestionId = c.Param("questionId")
		if a.HaveAnswer() && basic.MethodIsOk(c, "DELETE") && a.Delete() == nil {
			//
		}
		c.Abort()
	}
}

//查看回答中间件
func ViewAnswer() gin.HandlerFunc {
	return func(c *gin.Context) {
		a := features.NewAnswer()
		a.Uid = features.G_user.Info.Uid
		a.QuestionId = c.Param("questionId")
		if a.HaveAnswer() && basic.MethodIsOk(c, "GET") {
			a.Id = a.GetId()
			a.Content = a.GetContent()
			a.Time = a.GetTime()
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"data": gin.H{
					"uid":         a.Uid,
					"question_id": a.QuestionId,
					"answer_id":   a.Id,
					"time":        a.Time,
					"content":     a.Content,
				},
			})
		}
		c.Abort()
	}
}

//关注问题
func FollowQuestion() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := features.NewQuestion()
		q.Id = c.Param("questionId")
		if q.IsFollow() && basic.MethodIsOk(c, "DELETE") && q.CancelFollow() == nil {
			c.JSON(http.StatusOK, gin.H{
				"isFollow": "no",
			})
		} else if !q.IsFollow() && basic.MethodIsOk(c, "POST") && q.Follow() == nil {
			c.JSON(http.StatusOK, gin.H{
				"isFollow": "yes",
			})
		}
		c.Abort()
	}
}

//发表问题评论
func PostQuestionComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		qc := features.NewQuestionComment()
		qc.QuestionId = c.Param("questionId")
		qc.Time = basic.GetTimeNow()
		qc.Pid = c.DefaultPostForm("pid", "0")
		qc.Content = c.PostForm("content")
		qc.Id = basic.GetACommentId()
		qc.Uid = features.G_user.Info.Uid
		if basic.MethodIsOk(c, "POST") && qc.Post() == nil {
			//
		}
		c.Abort()
	}
}

//发表回答评论
func PostAnswerComments() gin.HandlerFunc {
	return func(c *gin.Context) {
		ac := features.NewAnswerComment()
		ac.AnswerId = c.Param("answerId")
		ac.Time = basic.GetTimeNow()
		ac.Pid = c.DefaultPostForm("pid", "0")
		ac.Content = c.PostForm("content")
		ac.Id = basic.GetACommentId()
		ac.Uid = features.G_user.Info.Uid
		if basic.MethodIsOk(c, "POST") && ac.Post() == nil {
			//
		}
		c.Abort()
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
