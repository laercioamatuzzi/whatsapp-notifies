package main

/* func main() {

	api := manager.Api{WhatsappWeb: connection.WhatsAppWeb{Number: "wpp-notify"}}
	api.Run()

}
*/

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"whatsapp-notifies/connection"
	"whatsapp-notifies/manager"

	"github.com/urfave/cli/v3"
)

func main() {

	api := manager.Api{WhatsappWeb: connection.WhatsAppWeb{Number: "wpp-notify"}}
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"cfg"},
				Usage:   "start the configuration process",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Hello, press enter to start the configuration process, after the api is running acess localhost:8080/qrcode to scan the QR code with your WhatsApp account...")
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()
					go api.Run()
					time.Sleep(time.Second * 5)
					for {
						fmt.Println("press enter after the QR code is scanned...")
						scanner = bufio.NewScanner(os.Stdin)
						scanner.Scan()
						if api.WhatsappWeb.IsConnected() {
							break
						}
						fmt.Println("Whatsapp is not connected, refresh the QR code and scan again...")
					}
					return nil
				},
			},
			{
				Name:    "start-app",
				Aliases: []string{"sa"},
				Usage:   "start application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					api.Run()
					return nil
				},
			},
			{
				Name:    "gateway",
				Aliases: []string{"gw"},
				Usage:   "start application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					gateway := manager.Gateway{}
					gateway.Init()
					fmt.Println(gateway.GetScheduleMessages())
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
