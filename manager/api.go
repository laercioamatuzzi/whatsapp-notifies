package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"whatsapp-notifies/connection"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Api struct {
	r           *gin.Engine
	WhatsappWeb connection.WhatsAppWeb
	SqliteConn  *connection.SqliteConn
}

func (a *Api) init(dbPath string) {
	a.SqliteConn = &connection.SqliteConn{DBPath: dbPath}
	a.r = gin.Default()
	a.WhatsappWeb = connection.WhatsAppWeb{Number: connection.WHATSAPPDB}
	a.WhatsappWeb.Login()

	a.r.Handle("GET", "ping", a.Ping)
	a.r.Handle("GET", "qrcode", a.GetQRCode)
	a.r.Handle("POST", "text", a.PostText)
	a.r.Handle("POST", "schedule", a.ScheduleMessage)
	a.r.Handle("GET", "schedule", a.GetScheduleMessages)
	a.r.Handle("GET", "schedule/all", a.GetAllMessages)

}

func (a *Api) Run(dbPath string) {
	a.init(dbPath)
	a.r.Run()
}

func (a *Api) Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})

}

func (a *Api) GetQRCode(c *gin.Context) {
	r := c.Request
	w := c.Writer

	browser := "Google Chrome (Linux)"
	version := [3]uint32{1, 0, 0}
	store.SetOSInfo(browser, version)
	store.DeviceProps.Os = &browser

	query := fmt.Sprintf("file:%s?_foreign_keys=on", connection.WHATSAPPDB)
	container, err := sqlstore.New("sqlite3", query, nil)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		go func() {
			for evt := range qrChan {
				if evt.Event == "code" {
					png, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
					w.Write(png)
					time.Sleep(time.Second * 30)
					line := connection.WhatsAppWeb{Number: connection.WHATSAPPDB}
					line.Login()

					if !line.IsConnected() {

						os.Remove(connection.WHATSAPPDB)
						return
					}

					a.WhatsappWeb = line

				} else {
					fmt.Println("Login event:", evt.Event)
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode("line already exist")
				}
			}
		}()
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("line already exist")
		return
	}

	time.Sleep(time.Second * 2)

}

func (a *Api) PostText(c *gin.Context) {

	// Parse JSON
	var json struct {
		Phone string `json:"phone" binding:"required"`
		Text  string `json:"text" binding:"required"`
	}

	err := c.Bind(&json)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	jid, err := a.WhatsappWeb.GetJid(json.Phone)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	messesageId, err := a.WhatsappWeb.SendText(jid, json.Text)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message_id": messesageId,
	})

}

func (a *Api) ScheduleMessage(c *gin.Context) {

	// Parse JSON
	var json struct {
		Phones []string `json:"phones" binding:"required"`
		Text   string   `json:"text" binding:"required"`
		Date   string   `json:"date" binding:"required"`
	}

	err := c.Bind(&json)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	jids := make([]types.JID, 0)
	for _, phone := range json.Phones {
		jid, err := a.WhatsappWeb.GetJid(phone)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		jids = append(jids, jid)

	}

	for _, jid := range jids {
		a.SqliteConn.Insert(jid.User, json.Text, json.Date)
	}

	c.JSON(200, gin.H{
		"message": "message scheduled",
	})

}

func (a *Api) GetScheduleMessages(c *gin.Context) {

	scheduleMessages := a.SqliteConn.GetWaitingScheduleMessages()

	c.JSON(200, gin.H{
		"message": scheduleMessages,
	})

}

func (a *Api) GetAllMessages(c *gin.Context) {

	scheduleMessages := a.SqliteConn.GetAllScheduleMessages()

	c.JSON(200, gin.H{
		"message": scheduleMessages,
	})

}
