package database

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

//数据库结构体，内置表格结构体，分别为数据库和结构体增加增删改查方法
type Database struct {
	UserName string //登录数据库的帐号和密码
	Password string
	DbName   string
	Db       *sql.DB //数据库
	table    Table
}

type Table struct {
	Db *sql.DB //表格所属的数据库
}

type AlterData struct {
	newKey     string
	newValue   string
	targeKey   string
	targeValue string
}
type FindData struct {
	targeKey   string
	targeValue string
}

//创建数据库
func (d *Database) Create(NewDbName string) error {
	stmt, err := d.Db.Prepare("create database " + NewDbName)
	stmt.Exec()
	return err
}

//打开数据库
func (d *Database) Open() error {
	command := strings.Join([]string{d.UserName, ":", d.Password, "@tcp(127.0.0.1:3306)/", d.DbName, "?charset=utf8"}, "")
	db, err := sql.Open("mysql", command)
	d.Db = db
	d.table.Db = db
	return err
}

//删除数据库
func (d *Database) Drop(DbName string) error {
	stmt, err := d.Db.Prepare("drop database " + DbName)
	stmt.Exec()
	return err
}

//关闭数据库
func (d *Database) Close() {
	d.Db.Close()
}

//创建表格
func (t *Table) Create(tableName string, kAndV string) error {
	command := strings.Join([]string{"create table ", tableName, " (id int NOT NULL auto_increment,", kAndV, ", primary key(id)) character set = utf8"}, "")
	stmt, err := t.Db.Prepare(command)
	stmt.Exec()
	return err
}

//删除表格
func (t *Table) Drop(tableName string) error {
	command := "drop table " + tableName
	stmt, err := t.Db.Prepare(command)
	stmt.Exec()
	return err
}

//插入记录
func (t *Table) Insert(tableName string, targeKey []string, targeValue []string) error {
	command := strings.Join([]string{"insert into ", tableName, "( ", strings.Join(targeKey, ""), " ) values", ObtainDbStr(targeValue)}, "")
	stmt, err := t.Db.Prepare(command)
	stmt.Exec(StrToInterface(targeValue)...)
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
func (t *Table) Alter(tableName string, data AlterData) error {
	command := strings.Join([]string{"update ", tableName, " set ", data.newKey, " = \"", data.newValue, "\" where ", data.targeKey + " = \"", data.targeValue, "\""}, "")
	stmt, err := t.Db.Prepare(command)
	stmt.Exec()
	return err
}

//查找记录
func (t *Table) Find(tableName string, limit string, targeKey string, targeValue string) ([]map[string]interface{}, error) {
	command := strings.Join([]string{"select ", limit, " from ", tableName, " where ", targeKey, " ='", targeValue, "'"}, "")
	stmt, err := t.Db.Query(command)
	columns, _ := stmt.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for stmt.Next() {
		_ = stmt.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	_ = stmt.Close()

	return list, err
}

type USER struct {
	Id           string
	Username     string
	Password     string
	Uid          string
	Gender       string
	Nickname     string
	Introduction string
	Avatar       string
	Question_id  string
	Reply_id     string
	Favorite_id  string
	Followers_id string
	Concern_id   string
	Article_id   string
}

func CheckError(err error, errorMsg string) {
	if err != nil {
		fmt.Println(errors.New(errorMsg))
	}
}
