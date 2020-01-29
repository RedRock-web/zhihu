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
		ViewMyInfo(targe, followingValus)
	} else {
		ViewOtherInfo(targe, followingValus)
	}
}

func ViewMyInfo(targe string, followingValus string) {
	ViewCommonInfo(targe, followingValus)
}

func ViewOtherInfo(targe string, followingValus string) {
	ViewCommonInfo(targe, followingValus)

}

func ViewCommonInfo(targe string, followingValus string) {
	switch targe {

	case "answers":
		{

		}
	case "asks":
		{

		}
	case "posts":
		{

		}
	case "columns":
		{

		}
	case "pins":
		{

		}
	case "collections":
		{

		}
	case "following":
		{

		}
	case "followers":
		{
			switch followingValus {
			case "columns":
				{

				}
			case "topics":
				{

				}
			case "questions":
				{

				}
			case "collections":
				{

				}

			}
		}
	}
}

func Chat(c *gin.Context)  {

}