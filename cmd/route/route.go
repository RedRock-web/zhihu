package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/features"
)

//使用route结构体，把站点各个页面模块作为方法
type Route struct {
	auth *gin.RouterGroup
}

func (route Route) Start() {
	r := gin.Default()
	LoginPage(r)
	//route.auth = r.Group("")
	//route.auth.Use(route.AuthRequired())
	//{
	//	route.HomePage()
	//	route.PersonalPage()
	//	route.QuestionDetailsPage()
	//}
	err := r.Run()
	basic.CheckError(err, "run失败！")
}

//通过cookie判断是否进入登录注册页
func LoginPageJudgeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if features.IsLogin(c, "userID") {
			basic.RediRect(c, "/")
			c.Next()
		} else {
			c.Next()
		}
	}
}

//登录注册页中间件
func LoginPageAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic.NowTimeUinx = basic.GetTimeUinx()
		c.SetCookie("login_page", basic.NowTimeUinx, 100, "/sign_in", "127.0.0.1", false, true)
		k, _ := c.Cookie("login_page")
		if k != "" {
			c.Next()
		} else {
			c.Abort()
		}
	}
}

//使用cookie检测是否登录的中间件
func (route Route) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("userID")
		if cookie == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请先登录！"})
			c.Abort()
		} else {
			c.Next()
			//c.JSON(200, "have cookie")
		}
	}
}

//登录注册页feature
func LoginPage(r *gin.Engine) {
	loginGroup := r.Group("")
	loginGroup.Use(LoginPageJudgeAuth())
	{
		r.POST("/sign_in", LoginPageJudgeAuth(), features.RegisteOrLogin)
		r.GET("/sign_in", LoginPageJudgeAuth())
		r.GET("/account/password_reset", LoginPageJudgeAuth(), features.PasswdReset)
	}
}

//主页
func (route Route) HomePage() {
	//主页
	route.auth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "这里是主页！",
		})
	})
	//搜索
	route.auth.GET("/search", )
	//推荐
	route.auth.GET("/feed/topstory/recommend", )
	//关注
	route.auth.GET("/feed/topstory/follow_wonderful", )
	//热榜
	route.auth.GET("/feed/topstory/hot", )
	//提问
	route.auth.POST("/questions", features.Quiz)
	//注销
	route.auth.GET("/logout", features.Logout)
}

//用户详情页
func (route Route) PersonalPage() {
	//编辑个人资料
	route.auth.GET("/edit", func(c *gin.Context) {
		route.auth.PUT("/me", features.Edit)
	})
	//查看用户状态
	route.auth.GET("/members/:nickname/:targe/:followingValus", )
	route.auth.POST("/chat", )
}

// 问题详情页
func (route Route) QuestionDetailsPage() {

	//进入问题详情页，获取问题信息
	route.auth.GET("/questions/:questionId/", func(c *gin.Context) {
		features.G_question_id = c.Param("questionId")

		//关注问题
		route.auth.POST("/questions/:questionId/followers", features.Follow)

		//取消关注问题
		route.auth.DELETE("/questions/:questionId/followers", features.CancelFollow)

		//查看评论
		route.auth.GET("/comments/:commentId", )

		//查看子评论
		route.auth.GET("/comments/:commentId/child_comments", )

		//对问题发表评论
		route.auth.POST("/questions/:questionId/comments", features.PostQuestionComments)

		//删除评论
		route.auth.DELETE("/comments/:commentId", features.DeleteComment)

		//点赞或反对评论
		route.auth.POST("comments/:commentId/actions/like")
		route.auth.DELETE("comments/:commentId/actions/like")
		route.auth.POST("comments/:commentId/actions/dislike")
		route.auth.DELETE("comments/:commentId/actions/dislike")

		//判断是否写了回答
		if features.HaveAnswer(features.G_question_id) {
			answerId := features.GetAnswerId(features.G_question_id)
			//查看自己的回答
			route.auth.GET("/questions/"+features.G_question_id+"/"+answerId, features.ViewAnswer)
			//删除自己的回答
			route.auth.DELETE("/answers/"+answerId, features.DeleteAnswer)
		} else {
			//写回答
			route.auth.POST("/questions/:questionId/draft", features.PostAnswer)
		}

		//对回答发表评论
		route.auth.POST("/questions/:questionId/answers/:answerId/comments", features.PostAnswerComments)

		//对回答的点赞
		route.auth.POST("/answers/:answerId/voters")
	})
}
