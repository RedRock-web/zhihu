package database

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
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
type InserData struct {
	username string
	password string
	uid      string
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
	command := d.UserName + ":" + d.Password + "@tcp(127.0.0.1:3306)/" + d.DbName + "?charset=utf8"
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
	command := "create table " + tableName + " (id int NOT NULL auto_increment," +
		kAndV + ", primary key(id)) character set = utf8"
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
//TODO:改为通用的insert
func (t *Table) Insert(tableName string, data InserData) error {
	command := "insert into " + tableName +
		"(username, password, uid) values(?,?,?)"
	stmt, err := t.Db.Prepare(command)
	stmt.Exec(data.username, data.password, data.uid)
	return err
}

//更改记录
func (t *Table) Alter(tableName string, data AlterData) error {
	command := "update " + tableName + " set " +
		data.newKey + " = \"" + data.newValue + "\" where " + data.targeKey + " = \"" + data.targeValue + "\""
	stmt, err := t.Db.Prepare(command)
	stmt.Exec()
	return err
}

//查找记录
func (t *Table) Find(tableName string, data FindData) (userImformation USER, err error) {
	command := "select * from " + tableName + " where " + data.targeKey + " = " + "\"" + data.targeValue + "\""
	stmt, err := t.Db.Query(command)
	for stmt.Next() {
		stmt.Scan(
			&userImformation.Id,
			&userImformation.Username,
			&userImformation.Password,
			&userImformation.Uid,
			&userImformation.Gender,
			&userImformation.Nickname,
			&userImformation.Introduction,
			&userImformation.Avatar,
			&userImformation.Question_id,
			&userImformation.Reply_id,
			&userImformation.Favorite_id,
			&userImformation.Followers_id,
			&userImformation.Concern_id,
			&userImformation.Article_id)
	}
	return
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
