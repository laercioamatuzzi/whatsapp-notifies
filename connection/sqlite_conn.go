package connection

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteConn struct {
	Database string
}

func (s *SqliteConn) Init() {

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table schedule (id integer not null primary key, phone text, text text, date datetime);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (s *SqliteConn) Insert(phone string, text string, date string) {

	db, err := sql.Open("sqlite3", "./whatsapp-notifies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	insert into schedule (phone, text, date) values (?, ?, ?);
	`
	_, err = db.Exec(sqlStmt, phone, text, date)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
