package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"whatsapp-notifies/utils"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteConn struct {
}

type status map[string]int

var statusMap = status{
	"scheduled": 0,
	"sent":      1,
	"delivered": 2,
	"read":      3,
	"failed":    -99,
}

var NOTIFIESDB = os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH") + utils.NOTIFIES_DB_NAME

func (s *SqliteConn) Migration() {

	db, err := sql.Open("sqlite3", NOTIFIESDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table schedule (id integer not null primary key, phone text, text text, date datetime, status integer);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (s *SqliteConn) Insert(phone string, text string, date string) {

	db, err := sql.Open("sqlite3", NOTIFIESDB)
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

	db, err := sql.Open("sqlite3", NOTIFIESDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := "select * from schedule where status = ? and date <= datetime('now')"
	rows, err := db.Query(sqlStmt, statusMap["scheduled"])

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

func (s *SqliteConn) UpdateMessageStatus(messageId int, statusTo int) {

	db, err := sql.Open("sqlite3", NOTIFIESDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	update schedule set status = ? where id = ?;
	`
	_, err = db.Exec(sqlStmt, statusTo, messageId)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
