# whatsapp-notifies

## Goal

This is a go application that allows you to schedule and send messages to your WhatsApp contacts.

## How to use?

### Prerequisites

- Go 1.23 or higher
- A WhatsApp account

### Installation

1. Clone the repository
2. Run `go build`

### Usage

#### Configuration

1. Run `whatsapp-notifies config --all`
2. Follow the instructions to scan the QR code with your WhatsApp account
3. After the QR code is scanned, it will setup the database and the configuration is done
4. Run `whatsapp-notifies start-app` to start the application

#### schedule a message using the command line

1. Run `whatsapp-notifies schedule --phone <PHONE> --text <TEXT> --date <DATE>`

#### schedule a message using the api

1. curl -X POST -H "Content-Type: application/json" -d '{"phone":"<PHONE>","text":"<TEXT>","date":"<DATE>"}' http://localhost:8080/schedule

#### get the schedule messages

1. curl http://localhost:8080/schedule 

