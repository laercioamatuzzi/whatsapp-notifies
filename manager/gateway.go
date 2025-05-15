package manager

import (
	"whatsapp-notifies/connection"
)

type Gateway struct {
	sqliteConn connection.SqliteConn
}

func (g *Gateway) Init() {
	g.sqliteConn = connection.SqliteConn{}

}

func (g *Gateway) GetScheduleMessages() []string {
	scheduleMessages := g.sqliteConn.GetScheduleMessages()
	return scheduleMessages
}
