package database

import (
	"database/sql"
	"zhihu/cmd/basic"
)

func CreateDatabase(db *sql.DB, NewDataBaseName string) {
	stmt, err := db.Prepare("create database " + NewDataBaseName)
	basic.CheckError(err)
	stmt.Exec()
}
func InsertField(db *sql.DB, tableName string, username string, password string) {
	stmt, err := db.Prepare("insert into " + tableName +
		"(username, password) values(?,?)")
	basic.CheckError(err)
	stmt.Exec(username, password)
}
func CreateTable(db *sql.DB, tableName string, keysAndValues string) {
	stmt, err := db.Prepare("create table " + tableName + " (id int NOT NULL auto_increment," +
		keysAndValues + ", primary key(id)) character set = utf8")
	basic.CheckError(err)
	stmt.Exec()
}

func OpenDatabase(username string, password string, databaseName string) (db *sql.DB) {
	db, err := sql.Open("mysql", username+":"+password+"@tcp(127.0.0.1:3306)/"+databaseName+"?charset=utf8")
	basic.CheckError(err)
	return db
}


func DatabaseSearch(db *sql.DB, tableName string, username string) (name string, passwd string) {
	var id string

	selectOder := "select * from " + tableName + " where username= \"" + username + "\""
	stmt, err := db.Query(selectOder)
	basic.CheckError(err)
	for stmt.Next() {
		stmt.Scan(&id, &name, &passwd)
	}
	return name, passwd
}

