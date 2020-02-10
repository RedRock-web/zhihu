package route

import (
	"github.com/gin-gonic/gin"
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
		loginPage.POST("/sign_in", RegisteOrLogin())
		loginPage.GET("/sign_in")
		loginPage.GET("/account/password_reset", )
	}
}

//用户详情页
func (e Engine) PersonalPage() {
	personalPage := e.r.Group("", Authorized2Some())
	{
		//编辑个人资料
		personalPage.GET("/edit")
		personalPage.PUT("/me", Edit())
		personalPage.POST("/chat", )
	}
}

//问题详情页
func (e Engine) QuestionPage() {
	e.Question()
	e.Answer()
}

//主页
func (e Engine) HomePage() {
	//主页
	e.r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "这里是主页！",
		})
	})
	//注销
	e.r.GET("/logout", Authorized2Some(), LogOut())
	//提问
	e.r.POST("/questions", Authorized2Some(), Quiz())

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

//问题详情页-回答
func (e Engine) Answer() {
	//需要登录
	RequiredLogin := e.r.Group("/answers", RefreshCookie(), Authorized2Some())
	{
		//对回答的表态
		RequiredLogin.POST("/:answerId/voters", VoteAnswer())

		//对回答发表评论
		RequiredLogin.POST("/:answerId/comments", PostAnswerComments())

		//删除回答评论
		RequiredLogin.DELETE("/:answerId/comments/:commentId", DeteleAnswerComment())

		//赞成或反对评论
		RequiredLogin.POST("/:answerId/comments/:commentId/voters", VoteComment())

	}
	//无需登录
	Nologin := e.r.Group("/answer")
	{
		//查看回答评论
		Nologin.GET("/:answerId/comments", ViewAnswerComment())

		//查看回答子评论
		Nologin.GET("/:answerId/child_comments/commentId", ViewChildAnswerComment())
	}
}

//问题详情页-问题
func (e Engine) Question() {
	//无需登录
	NoLogin := e.r.Group("/questions")
	{
		//进入问题详情页，获取问题信息
		NoLogin.GET("/:questionId/", GetQuestion())

		//查看问题评论
		NoLogin.GET("/:questionId/comments", ViewQuestionComment())

		//查看问题子评论
		NoLogin.GET("/:questionId/child_comments/:commentId", ViewChildQuestionComment())
	}

	//需登录
	q := e.r.Group("/questions", RefreshCookie(), Authorized2Some())
	{
		//关注问题
		q.POST("/:questionId/followers", FollowQuestion())

		//取消关注问题
		q.DELETE("/:questionId/followers", FollowQuestion())

		//对问题发表评论
		q.POST("/:questionId/comments", PostQuestionComments())

		//删除问题评论
		q.DELETE("/:questionId/comments/:commentId", DeteleQuestionComment())

		//点赞或反对评论
		q.POST("/:questionId/comments/:commentId/voters", VoteComment())

		//写回答
		q.POST("/:questionId/draft", ReplyAnswer())

		//查看自己的回答
		q.GET("/:questionId/answer", ViewAnswer())

		//删除自己的回答
		q.DELETE("/:questionId/answer", DeleteAnswer())
	}
}
