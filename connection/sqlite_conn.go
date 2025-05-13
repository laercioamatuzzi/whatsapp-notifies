package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	db, err := sql.Open("sqlite3", os.Getenv("SQLITE_DB_PATH"))
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

func (s *SqliteConn) GetScheduleMessages() []string {

	db, err := sql.Open("sqlite3", os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	select * from schedule;
	`
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	defer rows.Close()

	var scheduleMessages []string
	for rows.Next() {
		var phone string
		var text string
		var date string
		err := rows.Scan(&phone, &text, &date)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return nil
		}
		scheduleMessages = append(scheduleMessages, fmt.Sprintf("phone: %s, text: %s, date: %s", phone, text, date))
	}

	return scheduleMessages

}
