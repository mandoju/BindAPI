package Database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
)

var Db *sql.DB
var err error

func InitializeDb() {
	Db, err = sql.Open("mysql", "root:wp@tcp(18.220.253.183:8083)/bind")
	if err != nil {
		panic(err.Error())
	}
	file, err := ioutil.ReadFile("./query/createUsersTable.sql")
	if err != nil {
		panic(err.Error())
	}
	_, err = Db.Exec(string(file))
	if err != nil {
		panic(err.Error())
	}
}
