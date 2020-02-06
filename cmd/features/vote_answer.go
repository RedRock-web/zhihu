package features

import (
	"github.com/gin-gonic/gin"
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

func (av AnswerVote) Start(c *gin.Context) {
	targe := c.PostForm("type")

	if av.IsAgree() {
		av.Against() //如果已经点赞，那么无论点击赞同还是反对都是反对
	} else if av.IsAgainst() {
		av.Agree() //如果已经反对，那么无论点击反对还是赞同都是赞同
	} else { //没有表明态度
		if targe == "up" {
			av.Agree()
		} else if targe == "down" {
			av.Against()
		}
	}
}

//取消点赞
func (av AnswerVote) Against() {

}

//点赞
func (av AnswerVote) Agree() {

}

//判断是否反对
func (av AnswerVote) IsAgainst() bool {
	data, err := database.G_DB.Table.HighFind("answer_vote", "attitude", "answer_id = "+av.Id+"uid = "+av.Uid)
	basic.CheckError(err,"判断回答是否反对失败！")
	return string(data[0]["attitude"].([]uint8)) == "2"
}

//判断是否赞同
func (av AnswerVote) IsAgree() bool {
	data, err := database.G_DB.Table.HighFind("answer_vote", "attitude", "answer_id = "+av.Id+"uid = "+av.Uid)
	basic.CheckError(err,"判断回答是否反对失败！")
	return string(data[0]["attitude"].([]uint8)) == "1"
}
