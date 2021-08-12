package line

import (
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
)

var (
	lineSecretToken = os.Getenv("LINE_CHANNEL_SECRET")
	lineAccessToken = os.Getenv("LINE_CHANNEL_ACCESSTOKEN")
)

var (
	bot *linebot.Client
	err error
)

func init() {
	bot, err = linebot.New(lineSecretToken, lineAccessToken)
	if err != nil {
		logger.Error(err.Error())
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	events, cbErr := bot.ParseRequest(r)

	if cbErr != nil {
		if cbErr == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		logger.WithError(cbErr).Error("Line ParseRequest Error")
		return
	}

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			handleMessage(event)
		}
	}
}

func handleMessage(event *linebot.Event) {
	//var responseText string
	//var lineMsg []linebot.SendingMessage
	accountID, accountType := getAccountIDAndType(event)

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		//currently echo bot
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
			logger.WithField("func", "handleMessage").Error(err.Error())
			return
		}
		logger.WithFields(log.Fields{
			"func":        "handleMessage",
			"lineID":      accountID,
			"accountType": accountType,
		}).Infof("receive msg: %+v, ID: %+v", message.Text, message.ID)
		return
	}
}

const (
	accountTypeUser  = "user"
	accountTypeGroup = "group"
	accountTypeRoom  = "room"
)

func getAccountIDAndType(event *linebot.Event) (id, accountType string) {
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		return event.Source.UserID, accountTypeUser
	case linebot.EventSourceTypeGroup:
		return event.Source.GroupID, accountTypeGroup
	case linebot.EventSourceTypeRoom:
		return event.Source.RoomID, accountTypeRoom
	}
	return
}
