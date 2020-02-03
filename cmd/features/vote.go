package features

import (
	"github.com/gin-gonic/gin"
	"zhihu/cmd/database"
)

//Vote 表示一次用户对评论或回答的态度
type Vote struct {
	uid      string
	targeId  string
	attitude string
	time     string
}

//NewVote返回一次用户的表态
func NewVote() *Vote {
	return &Vote{}
}

func (v Vote) Start(c *gin.Context) {
	targe := c.PostForm("type")

	if v.IsAgree() {
		v.Against() //如果已经点赞，那么无论点击赞同还是反对都是反对
	} else if v.IsAgainst() {
		v.Agree() //如果已经反对，那么无论点击反对还是赞同都是赞同
	} else { //没有表明态度
		if targe == "up" {
			v.Agree()
		} else if targe == "down" {
			v.Against()
		}
	}
}

//取消点赞
func (v Vote) Against() {

}

//点赞
func (v Vote) Agree() {

}

//判断是否反对
func (v Vote) IsAgainst() bool {
	data, err := database.G_DB.Table.Find("", "id", "comment_id", c.id)

}

//判断是否赞同
func (v Vote) IsAgree() bool {

}
