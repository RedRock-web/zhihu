package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/internal/database"
)

//Info 表示用户信息
type Info struct {
	C            *gin.Context
	AlterTarge   string
	AlterContent string
	Uid          string //用户id
	Gender       string //性别
	Nickname     string //昵称
	Avatar       string //头像链接
	Introduction string
}

//更改个人信息
func (info Info) Alter() error {
	return database.G_DB.Table.Alter("user", info.AlterTarge, info.AlterContent, "uid", G_user.Info.Uid)
}

//查看个人信息
func (info *Info) View() error {
	info.Uid = G_user.Info.Uid

	nicknameTemp, err := database.G_DB.Table.Find("user", "nickname", "uid", G_user.Info.Uid)
	if err != nil {
		return err
	}
	info.Nickname = string(nicknameTemp[0]["nickname"].([]uint8))

	avatarTemp, _ := database.G_DB.Table.Find("user", "avatar", "uid", info.Uid)
	if avatarTemp[0]["avatar"] == nil {
	} else {
		info.Avatar = string(avatarTemp[0]["avatar"].([]uint8))
	}

	introductionTemp, _ := database.G_DB.Table.Find("user", "introduction", "uid", info.Uid)
	if introductionTemp[0]["avatar"] == nil {
	} else {
		info.Avatar = string(introductionTemp[0]["introduction"].([]uint8))
	}

	GenderTemp, err := database.G_DB.Table.Find("user", "gender", "uid", info.Uid)
	info.Gender = string(GenderTemp[0]["gender"].([]uint8))

	return err
}

//判断targe是否合理
func (info *Info) IsTargeCompliance() bool {
	return info.AlterTarge == "gender" || info.AlterTarge == "introduction" || info.AlterTarge == "nickname" || info.AlterTarge == "avatar"
}


