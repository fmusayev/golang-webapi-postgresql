package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // postgres driver
)

// DB ...
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

// ConnectSQL ...
func ConnectSQL(
	host string,
	port int,
	username string,
	password string,
	dbname string) (*DB, error) {

	connectionStringFormat := "host=%s port=%d user=%s password=%s dbname=%s"
	dataSource := fmt.Sprintf(connectionStringFormat, host, port, username, password, dbname)

	db, error := sql.Open("postgres", dataSource)
	if error != nil {
		log.Fatalln("error opening connection", error)
		panic(error)
	}
	dbConn.SQL = db
	return dbConn, error
}
