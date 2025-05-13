package connection

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/url"
	"os"
	"time"
	"whatsapp-notifies/utils"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

const (
	LOGLEVEL = "INFO"
)

type WhatsAppWeb struct {
	client *whatsmeow.Client
	Number string
}

func (w *WhatsAppWeb) eventHandler(evt interface{}) {

	switch v := evt.(type) {

	case *events.Message:

		if !v.Info.IsFromMe {
			from := v.Info.Sender.User
			messageId := v.Info.ID
			text := v.Message.GetConversation()
			image := v.Message.GetImageMessage()
			video := v.Message.GetVideoMessage()
			audio := v.Message.GetAudioMessage()
			document := v.Message.GetDocumentMessage()
			suggestion := v.Message.GetTemplateButtonReplyMessage()
			messageDate := time.Now().Format("2006-01-02 15:04:05")

			if len(text) <= 0 {
				extendedMsg := v.Message.GetExtendedTextMessage()
				text = extendedMsg.GetText()
			}

			var fileb64 string
			var mimeExt string
			var fileName string
			var suggestionText string
			var suggestionId string

			if image != nil {
				data, err := w.client.Download(image)
				if err != nil {
					utils.LoggingError(fmt.Sprintf("Failed to download image: %v", err))
					break
				}
				fileb64 = base64.StdEncoding.EncodeToString(data)
				ext, _ := mime.ExtensionsByType(image.GetMimetype())
				mimeExt = image.GetMimetype()

				if ext[0] == ".jpe" {
					fileName = fmt.Sprintf("%s.jpg", messageId)
				} else {
					fileName = fmt.Sprintf("%s%s", messageId, ext[0])
				}
			}

			if video != nil {
				data, err := w.client.Download(video)
				if err != nil {
					utils.LoggingError(fmt.Sprintf("Failed to download video: %v", err))
					break
				}
				fileb64 = base64.StdEncoding.EncodeToString(data)
				ext, _ := mime.ExtensionsByType(video.GetMimetype())
				mimeExt = video.GetMimetype()
				fileName = fmt.Sprintf("%s%s", messageId, ext[0])
			}

			if audio != nil {
				data, err := w.client.Download(audio)
				if err != nil {
					utils.LoggingError(fmt.Sprintf("Failed to download audio: %v", err))
					break
				}
				fileb64 = base64.StdEncoding.EncodeToString(data)
				ext, _ := mime.ExtensionsByType(audio.GetMimetype())
				mimeExt = audio.GetMimetype()
				fileName = fmt.Sprintf("%s%s", messageId, ext[0])
			}

			if document != nil {
				data, err := w.client.Download(document)
				if err != nil {
					utils.LoggingError(fmt.Sprintf("Failed to download document: %v", err))
					break
				}
				fileb64 = base64.StdEncoding.EncodeToString(data)
				ext, _ := mime.ExtensionsByType(document.GetMimetype())
				mimeExt = document.GetMimetype()
				fileName = fmt.Sprintf("%s%s", messageId, ext[0])
			}

			if suggestion != nil {
				suggestionId = suggestion.GetSelectedID()
				suggestionText = suggestion.GetSelectedDisplayText()
			}

			values := url.Values{
				"action":          {"message"},
				"sender":          {from},
				"receiver":        {w.GetClientNumber()},
				"message_id":      {messageId},
				"text":            {text},
				"file":            {fileb64},
				"file_type":       {mimeExt},
				"file_name":       {fileName},
				"suggestion_id":   {suggestionId},
				"suggestion_text": {suggestionText},
				"received_date":   {messageDate},
			}

			if (len(text) > 0 || len(fileb64) > 0) && !v.Info.IsGroup {
				fmt.Println(values)

			}

		}

	case *events.AppStateSyncComplete:
		if len(w.client.Store.PushName) > 0 && v.Name == appstate.WAPatchCriticalBlock {
			err := w.client.SendPresence(types.PresenceAvailable)
			if err != nil {
				utils.LoggingError(fmt.Sprintf("Failed to send available presence: %v", err))
			}
		}

	case *events.Connected, *events.PushNameSetting:
		if len(w.client.Store.PushName) == 0 {
			return
		}

		err := w.client.SendPresence(types.PresenceAvailable)
		if err != nil {
			utils.LoggingError(fmt.Sprintf("Failed to send available presence: %v", err))
		}

	case *events.Receipt:

		if !v.IsFromMe {
			eventType := "1"
			if v.Type == events.ReceiptTypeRead || v.Type == events.ReceiptTypeReadSelf {
				eventType = "3"
			} else if v.Type == events.ReceiptTypeDelivered {
				eventType = "2"
			}

			if eventType != "1" {
				for _, messageId := range v.MessageIDs {
					values := url.Values{
						"action":     {"dlr"},
						"sender":     {w.GetClientNumber()},
						"receiver":   {v.Sender.User},
						"message_id": {messageId},
						"type":       {eventType},
						"phone_name": {w.Number},
					}

					fmt.Println(values)

				}
			}
		}

	case *events.Contact:

		utils.LoggingInfo(fmt.Sprintf("CONTACT STUFF : %v", v))
		utils.LoggingInfo(fmt.Sprintf("CONTACT STUFF : %s, %s, %s", v.Action.GetFullName(), v.Action.GetFirstName(), v.JID.User))

	}
}

func (w *WhatsAppWeb) Login() {

	browser := "Google Chrome (Linux)"
	version := [3]uint32{1, 0, 0}
	store.SetOSInfo(browser, version)
	store.DeviceProps.Os = &browser

	dbLog := waLog.Stdout("Database", LOGLEVEL, true)
	query := fmt.Sprintf("file:%s?_foreign_keys=on", os.Getenv("WHATSAPP_DB_PATH"))
	container, err := sqlstore.New("sqlite3", query, dbLog)

	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()

	if err != nil {
		fmt.Printf("Failed to get device store: %v", err)
		return
	}
	clientLog := waLog.Stdout("Client", LOGLEVEL, true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(w.eventHandler)

	if client.Store.ID == nil {
		fmt.Printf("%s - SCAN QRCode is Required\n", w.Number)
		w.client = client
		return

	}

	err = client.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
		return
	}
	w.client = client
}

func (w *WhatsAppWeb) Disconnect() {
	w.client.Disconnect()
}

func (w *WhatsAppWeb) IsConnected() bool {

	return w.client.IsConnected()
}

func (w *WhatsAppWeb) GetClientNumber() string {

	return w.client.Store.ID.User
}

func (w *WhatsAppWeb) GetJid(number string) (types.JID, error) {

	if len(number) < 5 {
		return types.JID{}, fmt.Errorf("invalid number %s", number)
	}

	//w.client.IsConnected()
	isOnWhats, err := w.client.IsOnWhatsApp([]string{fmt.Sprintf("+%v", number)})

	if err != nil {
		utils.LoggingError(fmt.Sprintf("failed to check if number is on whatsapp. %v", err))
		return types.JID{}, err
	}

	return isOnWhats[0].JID, nil
}

func (w *WhatsAppWeb) SendText(jid types.JID, text string) (string, error) {

	messageId := w.client.GenerateMessageID()

	msg := &waProto.Message{Conversation: proto.String(text)}
	//ts, err := w.client.SendMessage(jid, messageId, msg)
	ts, err := w.client.SendMessage(context.Background(), jid, msg, whatsmeow.SendRequestExtra{ID: messageId})

	if err != nil {
		return "", err
	}

	utils.LoggingInfo(fmt.Sprintf("line [%s] Message Text: [%s] sent to: [%s] message_id: [%s] at: %s", w.Number, text, jid, messageId, ts.Timestamp))
	return messageId, nil
}
