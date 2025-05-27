package manager

import (
	"time"
	"whatsapp-notifies/connection"
)

type Gateway struct {
	sqliteConn  connection.SqliteConn
	WhatsappWeb connection.WhatsAppWeb
}

func (g *Gateway) Init() {
	g.sqliteConn = connection.SqliteConn{}

}

func (g *Gateway) GetScheduleMessages() []connection.ScheduleMessage {
	scheduleMessages := g.sqliteConn.GetScheduleMessages()
	return scheduleMessages
}

func (g *Gateway) SendMessage(scheduleMessage connection.ScheduleMessage) (string, error) {

	jid, err := g.WhatsappWeb.GetJid(scheduleMessage.Phone)
	if err != nil {
		return "", err
	}
	messesageId, err := g.WhatsappWeb.SendText(jid, scheduleMessage.Text)
	if err != nil {
		return "", err
	}

	return messesageId, nil

}

func (g *Gateway) Start() {

	g.Init()
	messagesWithError := []connection.ScheduleMessage{}
	messagesToRetry := []connection.ScheduleMessage{}

	for {

		scheduleMessages := g.GetScheduleMessages()
		for _, scheduleMessage := range scheduleMessages {
			messesageId, err := g.SendMessage(scheduleMessage)
			if err != nil {
				messagesToRetry = append(messagesToRetry, scheduleMessage)
				continue
			}

			err = g.sqliteConn.SetMessageId(messesageId, scheduleMessage.Id, connection.SENT)
			if err != nil {
				messagesWithError = append(messagesWithError, scheduleMessage)
			}

		}
		time.Sleep(time.Second * 2)
	}

}
