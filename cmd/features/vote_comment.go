package features

import (
	"zhihu/cmd/basic"
	"zhihu/cmd/database"
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

//点赞评论
func (cv CommentVote) Agree() error {
	err := database.G_DB.Table.Insert("comment_vote", []string{"uid", "time", "comment_id", "attitude"}, []string{cv.Uid, cv.Time, cv.Id, cv.Attitude})
	basic.CheckError(err, "点赞评论失败!")
	return err
}

//取消点赞评论
func (cv CommentVote) CancelAgree() error {
	err := database.G_DB.Table.Delete("comment_vote", "comment_id = "+cv.Id)
	basic.CheckError(err, "取消点赞失败!")
	return err
}

//反对评论
func (cv CommentVote) Against() error {
	err := database.G_DB.Table.Insert("comment_vote", []string{"uid", "time", "comment_id", "attitude"}, []string{cv.Uid, cv.Time, cv.Id, cv.Attitude})
	basic.CheckError(err, "点赞评论失败!")
	return err
}

//取消反对评论
func (cv CommentVote) CancelAgainst() error {
	err := database.G_DB.Table.Delete("comment_vote", "comment_id = "+cv.Id)
	basic.CheckError(err, "取消点赞失败!")
	return err
}

//获取态度,nil表不关心,1表赞同,0表反对
func (cv CommentVote) GetAttitude() string {
	data, err := database.G_DB.Table.HighFind("comment_vote", "attitude", "comment_id = "+cv.Id)
	basic.CheckError(err, "获取态度失败!")
		if data == nil {
			return ""
		} else {
			return string(data[0]["attitude"].([]uint8))
		}

}
