package database

import (
	"database/sql"
	"strings"
	"zhihu/cmd/basic"
)

var G_DB Database


//TODO:将time改为datatime格式
//项目数据库相关准备
func Start(){
	db := Database{
		UserName: "root",
		Password: "root",
		DbName:   "mysql",
	}
	err := db.Open()
	basic.CheckError(err, "打开数据库失败！")
	err = db.Create("zhihu")
	basic.CheckError(err, "创建数据库失败！")
	db1 := Database{
		UserName: "root",
		Password: "root",
		DbName:   "zhihu",
	}
	err = db1.Open()
	basic.CheckError(err, "打开数据库失败！")

	//user表格用于存储用户基本信息，包括
	//username    --  帐号
	//password    --  密码（加密后）
	//uid         --  用户id
	//gender      --  性别，0表示未知，1表示男，2表示女，默认0
	//nickname    --  用户名，默认‘知乎用户’
	//introdution --  个人介绍
	//avatar      --  头像链接
	err = db1.Table.Create("user", "username varchar(15), password varchar(40), uid int,gender char(2) not null DEFAULT '0',nickname varchar(20) not null DEFAULT '知乎用户', introduction varchar(200), avatar varchar(50)")
	basic.CheckError(err, "user表格创建失败！")

	//question表，用于存储问题基本信息
	//uid         --       提问的人
	//question_id --       问题id
	//time        --       问题创建时间
	//title       --       问题的标题
	//complement  --       问题的补充
	err = db1.Table.Create("question", "uid int, question_id int, time varchar(30), title varchar(30), complement varchar(300)")
	basic.CheckError(err, "question创建失败！")

	//answer表，用于存储回答基本信息
	//uid            --         回答的用户
	//question_id    --         回答的问题
	//answer_id      --         回答id
	//time           --         回答的时间
	//content        --         回答的内容
	err = db1.Table.Create("answer", "uid int, question_id int, answer_id int, time varchar(30),content varchar(4000)")
	basic.CheckError(err, "answer创建失败！")

	//question_comment表用于存储问题的评论
	//comment_id    --        评论的id
	//uid           --        评论的用户
	//question_id   --        评论的问题
	//time          --        评论的时间
	//content       --        评论的内容
	//pid存储comment_id，默认为0
	//若为0,则表示回复提问
	//若不为0,则表示向comment_id的评论回复
	err = db1.Table.Create("question_comment", "uid int, comment_id int, question_id int, pid int not null DEFAULT 0, time varchar(30), content varchar(200)")
	basic.CheckError(err, "question_comment创建失败！")

	//answer_comment表用于存储回答的评论
	//comment_id    --        评论的id
	//uid           --        评论的用户
	//answer_id     --        评论的回答
	//time          --        评论的时间
	//content       --        评论的内容
	//pid           --        评论的对象
	//pid存储表中的主键comment_id，默认为0
	//若为0,则表示回复回答
	//若不为0,则表示向comment_id评论的回复
	err = db1.Table.Create("answer_comment", "uid int, comment_id int, answer_id int, pid int not null DEFAULT 0, time varchar(30),content varchar(200)")
	basic.CheckError(err, "answer_comment创建失败！")

	//question_follow存用户对问题的关注
	//uid           --      用户
	//question_id   --		问题
	//如果存了某用户和某问题，则表示用户关注了此问题
	//如果没有存，则表示用户没有关注
	err = db1.Table.Create("question_follow", "uid int, question_id int")
	basic.CheckError(err, "question_follow表创建失败！")

	//user_follow存用户对其他的关注
	//uid           --      用户
	//follow_id     --      目标用户
	//如果存了某用户和目标用户，则表示用户关注了目标用户
	//如果没有存，则表示该用户没有关注目标用户
	err = db1.Table.Create("user_follow", "uid int, follow_uid int")
	basic.CheckError(err, "user_follow表创建失败！")

	//answer_vote存用户对答案的表态
	//uid           --      用户
	//time          --      时间
	//answer_id     --		回答
	//attitude      --      态度
	//没有表示不关心，1表示赞同，0表示反对
	err = db1.Table.Create("answer_vote", "uid int,time varchar(30), answer_id int, attitude int")
	basic.CheckError(err, "answer_vote表创建失败！")

	//comment_vote存用户对评论的态度
	//uid           --      用户
	//time          --      时间
	//comment_id    --		回答
	//attitude      --      态度
	//没有表示不关心，1表示赞同，0表示反对
	err = db1.Table.Create("comment_vote", "uid int, time varchar(30), comment_id int, attitude int")
	basic.CheckError(err, "comment_vote表创建失败！")

	//todo:增加收藏夹，文章，想法，专栏

	G_DB = db1
}

//数据库结构体，内置表格结构体，分别为数据库和结构体增加增删改查方法
type Database struct {
	UserName string //登录数据库的帐号和密码
	Password string
	DbName   string
	Db       *sql.DB //数据库
	Table    Table
}

type Table struct {
	Db *sql.DB //表格所属的数据库
}

//创建数据库
func (d *Database) Create(NewDbName string) error {
	stmt, err := d.Db.Prepare("create database " + NewDbName)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

//打开数据库
func (d *Database) Open() error {
	command := strings.Join([]string{d.UserName, ":", d.Password, "@tcp(127.0.0.1:3306)/", d.DbName, "?charset=utf8"}, "")
	db, err := sql.Open("mysql", command)
	d.Db = db
	d.Table.Db = db
	return err
}

//删除数据库
func (d *Database) Drop(DbName string) error {
	stmt, err := d.Db.Prepare("drop database " + DbName)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

//关闭数据库
func (d *Database) Close() error {
	return d.Db.Close()
}

//创建表格
func (t *Table) Create(tableName string, kAndV string) error {
	command := strings.Join([]string{"create table ", tableName, " (id int NOT NULL auto_increment,", kAndV, ", primary key(id)) character set = utf8"}, "")
	stmt, err := t.Db.Prepare(command)
	if err != nil {
		return err
	}
	_, err1 := stmt.Exec()
	return err1
}

//删除表格
func (t *Table) Drop(tableName string) error {
	command := "drop table " + tableName
	stmt, err := t.Db.Prepare(command)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

//插入记录
func (t *Table) Insert(tableName string, targeKey []string, targeValue []string) error {
	command := strings.Join([]string{"insert into ", tableName, "( ", strings.Join(targeKey, ","), " ) values ", ObtainDbStr(targeValue)}, "")
	//fmt.Println(command, targeValue)
	stmt, err := t.Db.Prepare(command)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(StrToInterface(targeValue)...)
	return err
}

//传入数组，返回一个特定string
func ObtainDbStr(data []string) string {
	str := "("
	for i := 0; i < len(data); i++ {
		str += "?,"
	}
	str = str[:len(str)-1]
	str += ")"
	return str
}

//字符串数组转空接口数组
func StrToInterface(data []string) []interface{} {
	data1 := make([]interface{}, len(data))
	for k, v := range data {
		data1[k] = v
	}
	return data1
}

//删除记录
func (t *Table) Delete(tableName string, limitInfo string) error {
	command := strings.Join([]string{"delete from ", tableName, "where ", limitInfo}, " ")
	//fmt.Println(command)
	stmt, err := t.Db.Prepare(command)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

//更改记录
func (t *Table) Alter(tableName string, newKey string, newValue string, limitKey string, limitValue string) error {
	command := strings.Join([]string{"update ", tableName, " set ", newKey, " = '", newValue, "' where ", limitKey + " = '", limitValue, "'"}, "")
	//fmt.Println(command)
	stmt, err := t.Db.Prepare(command)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

//查找记录,只能有一个限定条件
func (t *Table) Find(tableName string, limit string, targeKey string, targeValue string) ([]map[string]interface{}, error) {
	command := strings.Join([]string{"select ", limit, " from ", tableName, " where ", targeKey, " ='", targeValue, "'"}, "")
	//fmt.Println(command)
	stmt, err := t.Db.Query(command)
	if err != nil {
		return nil, err
	}
	columns, err := stmt.Columns()
	if err != nil {
		return nil, err
	}
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for stmt.Next() {
		err = stmt.Scan(cache...)
		if err != nil {
			return nil, err
		}
		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	err = stmt.Close()

	return list, err
}

//查找记录,可复杂限定条件
func (t *Table) HighFind(tableName string, targe string, limitInfo string) ([]map[string]interface{}, error) {
	command := strings.Join([]string{"select ", targe, " from ", tableName, " where ", limitInfo}, "")
	//fmt.Println(command)
	stmt, err := t.Db.Query(command)
	if err != nil {
		return nil, err
	}
	columns, err := stmt.Columns()
	if err != nil {
		return nil, err
	}
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for stmt.Next() {
		err = stmt.Scan(cache...)
		if err != nil {
			return nil, err
		}
		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	err = stmt.Close()

	return list, err
}

func UserName2Uid(username string) (string, error) {
	uid, err := G_DB.Table.Find("user", "uid", "username", username)
	return string(uid[0]["uid"].([]uint8)), err
}
func Uid2NickName(uid string) (string, error) {
	nickName, err := G_DB.Table.Find("user", "nickname", "uid", uid)
	return string(nickName[0]["nickname"].([]uint8)), err
}
