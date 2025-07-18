package connection

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SCHEDULED int = 0
	SENT      int = 1
	DELIVERED int = 2
	READ      int = 3
	FAILED    int = -99
)

type SqliteConn struct {
	DBPath string
}

type ScheduleMessage struct {
	Id     int    `json:"id"`
	Phone  string `json:"phone"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	Status int    `json:"status"`
}

func (s *SqliteConn) Migration() {

	db, err := sql.Open("sqlite3", s.DBPath)
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

	db, err := sql.Open("sqlite3", s.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	insert into schedule (phone, text, date, status) values (?, ?, ?, ?);
	`
	_, err = db.Exec(sqlStmt, phone, text, date, SCHEDULED)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (s *SqliteConn) GetScheduleMessages() []ScheduleMessage {

	db, err := sql.Open("sqlite3", s.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	select id, phone, text, date, status from schedule where status = ? and date <= datetime('now', 'localtime')
	`
	rows, err := db.Query(sqlStmt, SCHEDULED)

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
		scheduleMessages = append(scheduleMessages, ScheduleMessage{Id: id, Phone: phone, Text: text, Date: date, Status: status})
	}

	return scheduleMessages

}

func (s *SqliteConn) GetWaitingScheduleMessages() []ScheduleMessage {

	db, err := sql.Open("sqlite3", s.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	select id, phone, text, date, status from schedule where status = ?
	`
	rows, err := db.Query(sqlStmt, SCHEDULED)

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
		scheduleMessages = append(scheduleMessages, ScheduleMessage{Id: id, Phone: phone, Text: text, Date: date, Status: status})
	}

	return scheduleMessages

}

func (s *SqliteConn) GetAllScheduleMessages() []ScheduleMessage {

	db, err := sql.Open("sqlite3", s.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	select id, phone, text, date, status from schedule;
	`
	rows, err := db.Query(sqlStmt)

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
		scheduleMessages = append(scheduleMessages, ScheduleMessage{Id: id, Phone: phone, Text: text, Date: date, Status: status})
	}

	return scheduleMessages

}

func (s *SqliteConn) SetMessageId(messageId string, scheduleMessageId int, statusTo int) error {

	db, err := sql.Open("sqlite3", s.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	update schedule set message_id = ?, status = ? where id = ?;
	`
	_, err = db.Exec(sqlStmt, messageId, statusTo, scheduleMessageId)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}

func (s *SqliteConn) UpdateMessageStatus(messageId int, statusTo int) error {

	db, err := sql.Open("sqlite3", s.DBPath)
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
