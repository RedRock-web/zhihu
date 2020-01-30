package personal_page

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/basic"
)

type UserInfo struct {
}

//查看信息，分别查看个人和他人信息
func ViewInfo(c *gin.Context) {
	nickname := c.Param("nickname")
	targe := c.Param("targe")
	followingValus := c.Param("followingValus")

	if nickname == basic.G_NickName {
		ViewMyInfo(targe, followingValus, c)
	} else {
		ViewOtherInfo(targe, followingValus, c)
	}
}

func ViewMyInfo(targe string, followingValus string, c *gin.Context) {
	ViewCommonInfo(targe, followingValus, c)
}

func ViewOtherInfo(targe string, followingValus string, c *gin.Context) {
	ViewCommonInfo(targe, followingValus, c)
}

func ViewCommonInfo(targe string, followingValus string, c *gin.Context) {
	switch targe {

	case "answers": //回答
		{

		}
	case "asks": //提问
		{

		}
	case "posts": //文章
		{

		}
	case "columns": //专栏
		{

		}
	case "pins": //想法
		{

		}
	case "collections": //收藏
		{

		}
	case "followers": //关注者(粉丝)
		{

		}
	case "following": //关注了
		{
			switch followingValus {
			case "columns": //关注的专栏
				{

				}
			case "topics": //关注的话题
				{

				}
			case "questions": //关注的问题
				{

				}
			case "collections": //关注的收藏
				{

				}
			default: //关注的人

			}
		}
	default: //动态

	}
}

func Chat(c *gin.Context) {

}
