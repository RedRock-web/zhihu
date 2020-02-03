package features

//User 表示一个用户
type User struct {
	Account  Account  //帐号相关
	Info     Info     //用户信息
}

//NewUser 创建一个用户对象
func NewUser() *User {
	return &User{}
}


//当前用户信息缓存
var (
	G_user User
)

