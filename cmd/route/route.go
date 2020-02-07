package route

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/features"
)

//Engine 表示一个路由
type Engine struct {
	r *gin.Engine
}

//NewEngine 返回一个包含logger和recovery middleWare的路由
func NewEngine() *Engine {
	return &Engine{gin.Default()}
}

//main函数路由准备
func Start() {
	e := NewEngine()
	e.MiddleWare()
	e.Page()
	e.r.Run()
}

//设置全局路由middleWare
func (e Engine) MiddleWare() {
	e.r.Use(Cors())           //解决跨域问题
	e.r.NoRoute(NoResponse()) //解决404页面问题
}

//路由开启各项page feature
func (e Engine) Page() {
	e.LoginPage()
	e.PersonalPage()
	e.QuestionPage()
	e.HomePage()
}

//登录注册页feature
func (e Engine) LoginPage() {
	loginPage := e.r.Group("", Unauthorized2LoginPage())
	{
		loginPage.POST("/sign_in", features.RegisteOrLogin)
		loginPage.GET("/sign_in")
		loginPage.GET("/account/password_reset", features.PasswdReset)
	}
}

//用户详情页
func (e Engine) PersonalPage() {
	personalPage := e.r.Group("", Authorized2Some())
	{
		//编辑个人资料
		personalPage.GET("/edit")
		personalPage.PUT("/me", features.Edit)
		personalPage.POST("/chat", )
	}
}

//问题详情页
func (e Engine) QuestionPage() {
	e.Question()
	e.Answer()
	e.Comment()
}

//主页
func (e Engine) HomePage() {
	//主页
	e.r.GET("/", Authorized2Some(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "这里是主页！",
		})
	})
	//注销
	e.r.GET("/logout", Authorized2Some(), features.Logout)
	//提问
	e.r.POST("/questions", Authorized2Some(), features.Quiz)

	/*
		//搜索
		route.auth.GET("/search", )
		//推荐
		route.auth.GET("/feed/topstory/recommend", )
		//关注
		route.auth.GET("/feed/topstory/follow_wonderful", )
		//热榜
		route.auth.GET("/feed/topstory/hot", )
	*/

}

//问题详情页-评论
func (e Engine) Comment() {
	c := e.r.Group("/comments")
	{
		//删除评论
		c.DELETE("/:commentId", features.DeleteComment)
		//点赞或反对评论
		c.POST("/:commentId/voters")
	}
}

//问题详情页-回答
func (e Engine) Answer() {
	a := e.r.Group("/answers")
	{
		//对回答的点赞
		a.POST("/:answerId/voters")
		//对回答发表评论
		a.POST("/:answerId/comments", features.PostAnswerComments)
		//查看回答评论
		a.GET("/:answerId/comments")
		//查看回答子评论
		a.GET("/:answerId/child_comments/commentId", JudgeIfCommentInAnswer())
	}
}

//问题详情页-问题
func (e Engine) Question() {
	q := e.r.Group("/questions", RefreshCookie(),Authorized2Some())
	{
		//进入问题详情页，获取问题信息
		q.GET("/:questionId/")

		//关注问题
		q.POST("/:questionId/followers", JudgeIfFollowQuestion())

		//取消关注问题
		q.DELETE("/:questionId/followers", JudgeIfFollowQuestion())

		//对问题发表评论
		q.POST("/:questionId/comments", features.PostQuestionComments)
		//查看问题评论
		q.GET("/:questionId/comments")
		//查看问题子评论
		q.GET("/:questionId/child_comments/commentId", JudgeIfCommentInQuestion())
		//写回答
		q.POST("/:questionId/draft", JudgeIfReply(), features.PostAnswer)
		//查看自己的回答
		q.GET("/:questionId/answer", JudgeIfReply(), features.ViewAnswer)
		//删除自己的回答
		q.DELETE("/:questionId/answer", JudgeIfReply(), features.DeleteAnswer)
	}
}
