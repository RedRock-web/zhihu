package features

//Answer表示一个回答
type Answer struct {
	id          string //回答id
	time        string //回答时间
	uid         string //回答的用户
	question_id string //回答的问题
	content     string //回答的内容
}

//NewAnswer返回一个回答对象
func NewAnswer() *Answer {
	return &Answer{}
}
