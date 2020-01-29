package database

import (
	"database/sql"
	"strings"
	"zhihu/cmd/basic"
)

var G_DB Database

type USER struct {
	Id           string
	Username     string
	Password     string
	Uid          string
	Gender       string
	Nickname     string
	Introduction string
	Avatar       string
	QuestionId   string
	ReplyId      string
	FavoriteId   string
	FollowersId  string
	ConcernId    string
	ArticleId    string
}

//项目数据库相关准备
func Start() Database {
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
	err = db1.Table.Create("user", "username varchar(15), password varchar(15), uid int,gender int not null DEFAULT 0,nickname varchar(20) not null DEFAULT '知乎用户', introduction varchar(200), avatar varchar(50), question_id int, reply_id int, favorite_id int, followers_id int, concern_id int, article_id int")
	basic.CheckError(err, "user表格创建失败！")
	err = db1.Table.Create("comment", "uid int, targe int, targe_id int, pid int, time varchar(20)")
	basic.CheckError(err, "comment创建失败！")
	err = db1.Table.Create("question", "uid int, time int, title varchar(30), complement varchar(300), comment_id int, reply_id int")
	basic.CheckError(err, "question创建失败！")
	err = db1.Table.Create("reply", "uid int, question_id int, time varchar(20),comment_id int, attitude_reply_id int")
	basic.CheckError(err, "reply创建失败！")
	err = db1.Table.Create("reply_attitude", "reply_id int, uid int, attitude int")
	return db1
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

//更改记录
func (t *Table) Alter(tableName string, newKey string, newValue string, limitKey string, limitValue string) error {
	command := strings.Join([]string{"update ", tableName, " set ", newKey, " = '", newValue, "' where ", limitKey + " = '", limitValue, "'"}, "")
	stmt, err := t.Db.Prepare(command)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

//查找记录
func (t *Table) Find(tableName string, limit string, targeKey string, targeValue string) ([]map[string]interface{}, error) {
	command := strings.Join([]string{"select ", limit, " from ", tableName, " where ", targeKey, " ='", targeValue, "'"}, "")
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
	nickName, err := G_DB.Table.Find("user", "nickname", "uid",uid)
	return string(nickName[0]["nickname"].([]uint8)), err
}