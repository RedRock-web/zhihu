package features

import (
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
)

//AnswerVote 表示一次用户对回答的态度
type AnswerVote struct {
	Uid      string
	Id       string
	Attitude string
	Time     string
}

//NewAnswerVote返回一次用户对回答的表态
func NewAnswerVote() *AnswerVote {
	return &AnswerVote{}
}

//支持回答
func (av AnswerVote) Agree() error {
	err := database.G_DB.Table.Insert("answer_vote", []string{"uid", "time", "answer_id", "attitude"}, []string{av.Uid, av.Time, av.Id, av.Attitude})
	basic.CheckError(err, "支持回答失败!")
	return err
}

//取消支持回答
func (av AnswerVote) CancelAgree() error {
	err := database.G_DB.Table.Delete("answer_vote", "answer_id = "+av.Id)
	basic.CheckError(err, "取消支持回答失败!")
	return err
}

//反对回答
func (av AnswerVote) Against() error {
	err := database.G_DB.Table.Insert("answer_vote", []string{"uid", "time", "answer_id", "attitude"}, []string{av.Uid, av.Time, av.Id, av.Attitude})
	basic.CheckError(err, "反对回答失败!")
	return err
}

//取消反对回答
func (av AnswerVote) CancelAgainst() error {
	err := database.G_DB.Table.Delete("answer_vote", "answer_id = "+av.Id)
	basic.CheckError(err, "取消反对回答失败!")
	return err
}

//获取态度,nil表不关心,1表赞同,0表反对
func (av AnswerVote) GetAttitude() string {
	data, err := database.G_DB.Table.HighFind("answer_vote", "attitude", "answer_id = "+av.Id)
	basic.CheckError(err, "获取回答态度失败!")
	if data == nil {
		return ""
	} else {
		return string(data[0]["attitude"].([]uint8))
	}
}
