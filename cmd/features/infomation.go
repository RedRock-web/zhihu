package features

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//Info 表示用户信息
type Info struct {
	Uid      string //用户id
	Gender   string //性别
	Nickname string //昵称
	Avatar   string //头像链接
}

//更改个人信息接口
func Edit(c *gin.Context) {
	targe := c.PostForm("targe")
	content := c.PostForm("content")
	G_user.Info.Alter(c, targe, content)
}

//更改个人信息
func (info Info) Alter(c *gin.Context, targe string, content string) (error){
	if IsTargeCompliance(targe) {
		err := database.G_DB.Table.Alter("user", targe, content, "uid", G_user.Info.Uid)
		basic.CheckError(err, "更改个人信息失败！")
		if err == nil {
			err = info.View()
			basic.CheckError(err, "查看个人信息失败！")
			if err != nil {
				return err
			}
			c.JSON(http.StatusOK, gin.H{
				"uid":      info.Uid,
				"nickname": info.Nickname,
				"gender":   info.Gender,
				"avatar":   info.Avatar,
			})
		} else {
			return err
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"message": "非法操作！",
			},
		})
	}
	return nil
}

func (info Info) View() error {
	info.Uid = G_user.Info.Uid
	nicknameTemp, err := database.G_DB.Table.Find("user", "nickname", "uid", G_user.Info.Uid)
	if err != nil {
		return err
	}
	info.Nickname = string(nicknameTemp[0]["nickname"].([]uint8))
	avatarTemp, err := database.G_DB.Table.Find("user", "abatar", "uid", info.Uid)
	if err != nil {
		return err
	}
	info.Avatar = string(avatarTemp[0]["avatar"].([]uint8))
	GenderTemp, err := database.G_DB.Table.Find("user", "gender", "uid", info.Uid)
	info.Gender = string(GenderTemp[0]["avatar"].([]uint8))

	return err
}

//判断targe是否合理
func IsTargeCompliance(targe string) bool {
	return targe == "gender" || targe == "information" || targe == "nickname" || targe == "avatar"
}

//查看信息，分别查看个人和他人信息
func ViewInfo(c *gin.Context) {
	nickname := c.Param("nickname")
	targe := c.Param("targe")
	followingValus := c.Param("followingValus")

	if nickname == G_user.Info.Nickname {
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
