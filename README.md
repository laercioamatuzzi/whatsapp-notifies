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

1. Run `whatsapp-notifies config`
2. Follow the instructions to scan the QR code with your WhatsApp account
3. After the QR code is scanned, it will setup the database and the configuration is done
4. Run `whatsapp-notifies start-app` to start the application

#### schedule a message using the command line

1. Run `whatsapp-notifies schedule --phone <PHONE> --text <TEXT> --date <DATE>`

#### schedule a message using the api

1. curl -X POST -H "Content-Type: application/json" -d '{"phone":"<PHONE>","text":"<TEXT>","date":"<DATE>"}' http://localhost:8080/schedule

#### get the schedule messages

1. curl http://localhost:8080/schedule 


### Endpoints

#### ping - just to check if the api is running

```
curl http://localhost:8080/ping
```

#### qrcode - get the QR code to scan with your WhatsApp account

```
curl http://localhost:8080/qrcode
```

#### text - send a text message with the linked WhatsApp account to a phone number

```
curl -X POST -H "Content-Type: application/json" -d '{"phone":"<PHONE>","text":"<TEXT>"}' http://localhost:8080/text
```

#### schedule - schedule a message to send with the linked WhatsApp account

```
curl -X POST -H "Content-Type: application/json" -d '{"phones":["<PHONE1>","<PHONE2>"],"text":"<TEXT>","date":"<DATE>"}' http://localhost:8080/schedule
```

#### get schedule messages - list all the scheduled messages that will be sent by the gateway

```
curl http://localhost:8080/schedule 
