package connection

import (
	"database/sql"
	"log"
	"os"
	"whatsapp-notifies/utils"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SCHEDULE  = 0
	SENT      = 1
	DELIVERED = 2
	READ      = 3
	FAILED    = -99
)

type SqliteConn struct {
}

type ScheduleMessage struct {
	Id     int    `json:"id"`
	Phone  string `json:"phone"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	Status int    `json:"status"`
}

var NOTIFIESDB = os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH") + utils.NOTIFIES_DB_NAME

func (s *SqliteConn) Migration() {

	db, err := sql.Open("sqlite3", NOTIFIESDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table schedule (id integer not null primary key, phone text, text text, date datetime, status integer, message_id text default null);
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
	insert into schedule (phone, text, date, status) values (?, ?, ?, ?);
	`
	_, err = db.Exec(sqlStmt, phone, text, date, SCHEDULE)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (s *SqliteConn) GetScheduleMessages() []ScheduleMessage {

	db, err := sql.Open("sqlite3", NOTIFIESDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	select * from schedule where status = ? and date <= datetime('now')
	`
	rows, err := db.Query(sqlStmt, SCHEDULE)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	defer rows.Close()

	var scheduleMessages []ScheduleMessage
	for rows.Next() {
		var id int
		var phone string
		var text string
		var date string
		var status int
		err := rows.Scan(&id, &phone, &text, &date, &status)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return nil
		}
		scheduleMessages = append(scheduleMessages, ScheduleMessage{Phone: phone, Text: text, Date: date, Status: status})
	}

	return scheduleMessages

}

func (s *SqliteConn) SetMessageId(messageId string, scheduleMessageId int, statusTo int) error {

	db, err := sql.Open("sqlite3", NOTIFIESDB)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	update schedule set message_id = ? and status = ? where id = ?;
	`
	_, err = db.Exec(sqlStmt, messageId, statusTo, scheduleMessageId)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (s *SqliteConn) UpdateMessageStatus(messageId int, statusTo int) error {

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
		return err
	}
	return nil
}
