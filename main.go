package main

import (
	"whatsapp-notifies/connection"
	"whatsapp-notifies/manager"
)

func main() {

	api := manager.Api{WhatsappWeb: connection.WhatsAppWeb{Number: "wpp-notify"}}
	api.Run()

}
