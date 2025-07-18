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

### Version
    - 0.0.1

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

#### ping 

##### a simple ping to check if the application is running

```
curl http://localhost:8080/ping
```

#### qrcode 

##### get the QR code to scan with your WhatsApp account and link with the application

```
curl http://localhost:8080/qrcode
```

#### text 

##### send a text message to a WhatsApp contact using the linked WhatsApp account

```
curl -X POST -H "Content-Type: application/json" -d '{"phone":"<PHONE>","text":"<TEXT>"}' http://localhost:8080/text
```

#### schedule 

##### schedule a message to be sent by the gateway process using the linked WhatsApp account

```
curl -X POST -H "Content-Type: application/json" -d '{"phones":["<PHONE1>","<PHONE2>"],"text":"<TEXT>","date":"<DATE>"}' http://localhost:8080/schedule
```

#### get schedule messages 

##### get the scheduled messages to be sent by the gateway process

```
curl http://localhost:8080/schedule 
```

### TODO's

- [ ] Store whatsapp events
- [ ] Store whatsapp received messages
- [ ] Add groups support
- [ ] Endpoint to list received messages
- [ ] Endpoint to list events
- [ ] Frontend
- [ ] Finish Gateway process
- [ ] Add unit tests