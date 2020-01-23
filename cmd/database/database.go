package database

import (
	"database/sql"
	"fmt"
	"zhihu/cmd/basic"
)

func DatabasePrepare() (db *sql.DB) {
	db1 := OpenDatabase("root", "root", "mysql")
	defer db1.Close()
	//zhihu即项目数据库
	CreateDatabase(db1, "zhihu")
	db = OpenDatabase("root", "root", "zhihu")
	//defer db.Close()
	//user_information存放用户基本信息
	//TODO:字段推荐不使用NULL， 但是不使用NULL又不设置默认值，后续都无法对数据库操作，需要改进
	//gender 0代表女， 1代表男
	CreateTable(db, "user_information", "username varchar(20), password varchar(20),gender varchar(2), nickname varchar(20), introduction varchar(200), avatar varchar(40), question")

	return db
}

func CreateDatabase(db *sql.DB, NewDataBaseName string) error {
	stmt, err := db.Prepare("create database " + NewDataBaseName)
	basic.CheckError(err, "数据库创建失败！")
	stmt.Exec()
	return err
}
func InsertField(db *sql.DB, tableName string, username string, password string) error {
	stmt, err := db.Prepare("insert into " + tableName +
		"(username, password) values(?,?)")
	fmt.Println("insert into " + tableName +
		"(username, password) values(?,?)")
	basic.CheckError(err, "数据库插入出错！")
	stmt.Exec(username, password)
	return err
}
func CreateTable(db *sql.DB, tableName string, keysAndValues string) error {
	stmt, err := db.Prepare("create table " + tableName + " (id int NOT NULL auto_increment," +
		keysAndValues + ", primary key(id)) character set = utf8")
	basic.CheckError(err, "数据库表格创建失败！")
	stmt.Exec()
	return err
}

func OpenDatabase(dbUsername string, dbPassword string, databaseName string) (db *sql.DB) {
	db, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@tcp(127.0.0.1:3306)/"+databaseName+"?charset=utf8")
	basic.CheckError(err, "打开数据库失败！")
	return
}

func DatabaseSearch(db *sql.DB, tableName string, targeKey string, targeValue string) (userImformation basic.USER, err error) {
	selectOder := "select * from " + tableName + " where " + targeKey + " = " + "\"" + targeValue + "\""
	stmt, err := db.Query(selectOder)
	basic.CheckError(err, "数据库查找失败！")
	for stmt.Next() {
		stmt.Scan(&userImformation.Id, &userImformation.Username, &userImformation.Password, &userImformation.Gender, &userImformation.Nickname, &userImformation.Introduction, &userImformation.Avatar)
	}
	return userImformation, err
}

func DatabaseUpdate(db *sql.DB, tableName string, targeKey string, targeValue string, newKey string, newValue string) (error){
	oder := "update " + tableName + " set " +
		newKey + " = \"" + newValue + "\" where " + targeKey + " = \"" + targeValue + "\""
	fmt.Println(oder)
	stmt, err := db.Prepare(oder)
	basic.CheckError(err, "数据库修改失败！")
	stmt.Exec()
	return err
}
