package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL Driver
)

//DemoDB :
var DemoDB *sql.DB

//ConnectDatabase : Connecting database
func ConnectDatabase() (err error) {
	user := "root"
	password := ""
	DemoDB, err = sql.Open("mysql", user+":"+password+"@tcp(localhost:3306)/demo")
	if err != nil {
		return errors.New("DemoDB error :" + err.Error())
	}
	if err = DemoDB.Ping(); err != nil {
		return errors.New("DemoDB error :" + err.Error())
	}
	fmt.Println("database is connected.")

	return nil
}
