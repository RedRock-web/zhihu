package features

import (
	"github.com/gin-gonic/gin"
)

//CommentVote 表示一次用户对评论的态度
type CommentVote struct {
	Uid      string
	Id       string
	Attitude string
	Time     string
}

//NewCommentVote返回一次用户对评论的表态
func NewCommentVote() *CommentVote {
	return &CommentVote{}
}

func (cv CommentVote) Start(c *gin.Context) {
	targe := c.PostForm("type")

	if cv.IsAgree() {
		cv.Against() //如果已经点赞，那么无论点击赞同还是反对都是反对
	} else if cv.IsAgainst() {
		cv.Agree() //如果已经反对，那么无论点击反对还是赞同都是赞同
	} else { //没有表明态度
		if targe == "up" {
			cv.Agree()
		} else if targe == "down" {
			cv.Against()
		}
	}
}

//取消点赞
func (cv CommentVote) Against() {

}

//点赞
func (cv CommentVote) Agree() {

}

//判断是否反对
func (cv CommentVote) IsAgainst() bool {
	//data, err := database.G_DB.Table.Find("", "id", "comment_id", "")
	return true
}

//判断是否赞同
func (cv CommentVote) IsAgree() bool {
	return true
}
