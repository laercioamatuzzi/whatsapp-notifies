package manager

import (
	"os"
	"whatsapp-notifies/connection"
)

type Gateway struct {
	sqliteConn connection.SqliteConn
}

func (g *Gateway) Init() {
	sqliteConn := connection.SqliteConn{Database: os.Getenv("SQLITE_DB_PATH")}
	g.sqliteConn = sqliteConn

}

func (g *Gateway) GetScheduleMessages() []string {
	scheduleMessages := g.sqliteConn.GetScheduleMessages()
	return scheduleMessages
}
