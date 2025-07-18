package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"whatsapp-notifies/manager"
	"whatsapp-notifies/utils"

	"github.com/urfave/cli/v3"
)

const (
	VERSION = "0.0.1"
)

func main() {

	utils.LoadEnv()
	api := manager.Api{}
	gateway := manager.Gateway{WhatsappWeb: &api.WhatsappWeb}
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "./whatsapp-notifies config",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Hello, press enter to start the configuration process, after the api is running access http://localhost:8080/qrcode to scan the QR code with your WhatsApp account...")
					scanner := bufio.NewScanner(os.Stdin)
					scanner.Scan()
					go api.Run(os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH") + utils.NOTIFIES_DB_NAME)
					time.Sleep(time.Second * 2)
					for {
						fmt.Println("Acess http://localhost:8080/qrcode scan the QR code with your WhatsApp and wait for the Login event: `[Client INFO] Successfully authenticated` message in the terminal, then pressn enter to finish.")
						scanner = bufio.NewScanner(os.Stdin)
						scanner.Scan()
						if api.WhatsappWeb.IsConnected() {
							api.SqliteConn.Migration()
							break
						}
						fmt.Println("Whatsapp is not connected, refresh the QR code and scan again...")
					}
					return nil
				},
			},
			{
				Name:  "start-app",
				Usage: "start application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					go gateway.Start(os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH") + utils.NOTIFIES_DB_NAME)
					api.Run(os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH") + utils.NOTIFIES_DB_NAME)
					return nil
				},
			},
			{
				Name:    "gateway",
				Aliases: []string{"-g"},
				Usage:   "start application",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					gateway.Init(os.Getenv("WHATSAPP_NOTIFIES_CONFIG_PATH") + utils.NOTIFIES_DB_NAME)
					fmt.Println(gateway.GetScheduleMessages())
					return nil
				},
			},
			{
				Name:      "schedule",
				Usage:     "whatsapp-notifies schedule phone 35199999999 text hello date '2025-01-01 12:00:00'",
				ArgsUsage: "phone <PHONE> text <TEXT> date <DATE>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					api.SqliteConn.Insert(cmd.Args().Get(0), cmd.Args().Get(1), cmd.Args().Get(2))
					return nil
				},
			},
			{
				Name:      "version",
				Usage:     "whatsapp-notifies version",
				ArgsUsage: "",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("whatsapp-notifies version ", VERSION)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
