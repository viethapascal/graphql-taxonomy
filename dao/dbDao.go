package dao

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DbDao struct {
	Host string
	Port string
	User string
	Password string
	DbName string
}
var Db *sql.DB

func init()  {
	log.Println("init dao")
	d := DbDao{
		User:"admin",
		DbName: "newdb",
	}
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", d.User, d.DbName)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	Db = conn
}
