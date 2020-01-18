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
	CreateTable(db, "user_information", "username varchar(20), password varchar(20) ,gender varchar(2), nickname varchar(2), introduction varchar(10), avatar varchar(10)")

	return db
}

func CreateDatabase(db *sql.DB, NewDataBaseName string) {
	stmt, err := db.Prepare("create database " + NewDataBaseName)
	basic.CheckError(err, "数据库创建失败！")
	stmt.Exec()
}
func InsertField(db *sql.DB, tableName string, username string, password string) {
	stmt, err := db.Prepare("insert into " + tableName +
		"(username, password) values(?,?)")
	basic.CheckError(err, "数据库插入出错！")
	stmt.Exec(username, password)
}
func CreateTable(db *sql.DB, tableName string, keysAndValues string) {
	stmt, err := db.Prepare("create table " + tableName + " (id int NOT NULL auto_increment," +
		keysAndValues + ", primary key(id)) character set = utf8")
	basic.CheckError(err, "数据库表格创建失败！")
	stmt.Exec()
}

func OpenDatabase(dbUsername string, dbPassword string, databaseName string) (db *sql.DB) {
	db, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@tcp(127.0.0.1:3306)/"+databaseName+"?charset=utf8")
	basic.CheckError(err, "打开数据库失败！")
	return
}

func DatabaseSearch(db *sql.DB, tableName string, targeKey string, targeValue string) (userImformation basic.USER) {
	selectOder := "select * from " + tableName + " where " + targeKey + " = " + "\"" + targeValue + "\""
	stmt, err := db.Query(selectOder)
	basic.CheckError(err, "数据库查找失败！")
	for stmt.Next() {
		stmt.Scan(&userImformation.Id, &userImformation.Username, &userImformation.Password, &userImformation.Gender, &userImformation.Nickname, &userImformation.Introduction, &userImformation.Avatar)
	}
	return
}

func DatabaseUpdate(db *sql.DB, tableName string, targeKey string, targeValue string, newKey string, newValue string) {
	oder := "update " + tableName + " set " +
		newKey + " = \"" + newValue + "\" where " + targeKey + " = " + targeValue
	fmt.Println(oder)
	stmt, err := db.Prepare(oder)
	basic.CheckError(err, "数据库修改失败！")
	stmt.Exec()
}
