package manager

import (
	"fmt"
	"time"
	"whatsapp-notifies/connection"
	"whatsapp-notifies/utils"
)

type Gateway struct {
	sqliteConn  *connection.SqliteConn
	WhatsappWeb *connection.WhatsAppWeb
}

func (g *Gateway) Init(dbPath string) {
	g.sqliteConn = &connection.SqliteConn{DBPath: dbPath}

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

func (g *Gateway) Start(dbPath string) {

	g.Init(dbPath)
	//messagesWithError := []connection.ScheduleMessage{}
	//messagesToRetry := []connection.ScheduleMessage{}

	for {

		scheduleMessages := g.GetScheduleMessages()
		for _, scheduleMessage := range scheduleMessages {

			messageID, err := g.SendMessage(scheduleMessage)
			if err != nil {
				utils.LoggingError(fmt.Sprintf("Failed to send message: %v", err))
				//messagesToRetry = append(messagesToRetry, scheduleMessage)
				continue
			}

			err = g.sqliteConn.SetMessageId(messageID, scheduleMessage.Id, connection.SENT)
			if err != nil {
				//messagesWithError = append(messagesWithError, scheduleMessage)
				utils.LoggingError(fmt.Sprintf("Failed to set message id: %v", err))
			}

		}
		time.Sleep(time.Second * 5)
	}

}
