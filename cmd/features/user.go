package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//User 表示一个用户
type User struct {
	Account Account //帐号相关
	Info    Info    //用户信息
}

//NewUser 创建一个用户对象
func NewUser() *User {
	return &User{}
}

//当前用户信息缓存
var (
	G_user User
)

//获取我关注的人
func (u User) GetFollowing(limitInfo string) ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.HighFind("user_follow", "uid", "follow_uid = "+limitInfo)
	return data, err
}

//获取关注我的人
func (u User) GetFollowers(limitInfo string) ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.HighFind("user_follow", "follow_uid", "uid = "+limitInfo)
	return data, err
}

//发送json用户信息
func PostUser(c *gin.Context, u *User, uids []map[string]interface{}) {
	var user []gin.H

	for _, v := range uids {
		if v["uid"] == nil {
			u.Info.Uid = string(v["follow_uid"].([]uint8))
		} else {
			u.Info.Uid = string(v["uid"].([]uint8))
		}
		data, err := u.GetInfo()
		basic.CheckError(err, "获取用户信息失败!")
		user = append(user, gin.H{
			"uid":          u.Info.Uid,
			"gender":       string(data[0]["gender"].([]uint8)),
			"nickname":     string(data[0]["nickname"].([]uint8)),
			"introduction": string(data[0]["introduction"].([]uint8)),
			"avatar":       string(data[0]["avatar"].([]uint8)),
		})
	}
	//所有问题组合后,返回json
	c.JSON(200, gin.H{
		"status": 0,
		"data": gin.H{
			"user": user,
		},
	})
	c.Abort()
}

//获取用户信息
func (u User) GetInfo() ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.HighFind("user", "gender, nickname, introduction,avatar", "uid = "+u.Info.Uid)
	return data, err
}

//获取用户的提问
func (u User) GetQuestion() ([]map[string]interface{}, error) {
	data, err := database.G_DB.Table.HighFind("question", "question_id, time, title, complement", "uid = "+u.Info.Uid)
	return data, err
}
