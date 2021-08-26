package line

import (
	"fmt"
	"net/http"
	"os"

	"github.com/m1stborn/mistChatbot/internal/pkg/command"
	"github.com/m1stborn/mistChatbot/internal/pkg/model"

	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/sirupsen/logrus"
)

var (
	lineSecretToken = os.Getenv("LINE_CHANNEL_SECRET")
	lineAccessToken = os.Getenv("LINE_CHANNEL_ACCESSTOKEN")
	callbackUrl     = os.Getenv("CALLBACK_URL")
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
		case linebot.EventTypeFollow:
			handleFollow(event)
		case linebot.EventTypeUnfollow:
			handleUnfollow(event)
		}

	}
}

func handleMessage(event *linebot.Event) {
	var responseText string
	//var lineMsg []linebot.SendingMessage
	accountID, accountType := getAccountIDAndType(event)

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		//TODO handleFollow to ensure user exist in database
		user := model.DB.GetUser(accountID)
		//if !model.DB.CheckLineAccessTokenExist(accountID) {
		if user.LineAccessToken == "" {
			//TODO currently echo with not register user
			//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not connect with line notify")).Do() err != nil {
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
				logger.WithField("func", "handleMessage").Error(err.Error())
			}
			return
		}

		responseText = command.HandleCommand(message.Text, user, true)

		if responseText != "" {
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(responseText)).Do(); err != nil {
				logger.WithField("func", "handleMessage").Error(err.Error())
			}
		}

		logger.WithFields(log.Fields{
			"func":        "handleMessage",
			"lineID":      accountID,
			"accountType": accountType,
		}).Infof("receive msg: %+v, ID: %+v", message.Text, message.ID)
		return
	}
}

func handleUnfollow(event *linebot.Event) {
	accountID, _ := getAccountIDAndType(event)

	model.DB.UserUnfollow(accountID)
	model.DB.DeleteSubUserUnfollow(accountID)

	logger.WithFields(log.Fields{
		"func":   "HandleUnfollow",
		"lineID": accountID,
	}).Info("Line unfollow")
}

func handleFollow(event *linebot.Event) {
	accountID, _ := getAccountIDAndType(event)

	user := model.DB.GetUser(accountID)
	if user == nil {
		model.DB.CreateUser(&model.User{
			Line: accountID,
		})
	}

	url := getAuthorizeURL(accountID)
	//text := fmt.Sprintf("請至以下網址連動LINE NOTIFY與mistChatbot:\n,%s", callbackUrl+"/line/notify/auth")
	text := fmt.Sprintf("請至以下網址連動LINE NOTIFY與mistChatbot:\n,%s", url)
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		logger.WithFields(log.Fields{
			"func":   "handleFollow",
			"lineID": accountID,
		}).Error(err)
		return
	}
	logger.WithFields(log.Fields{
		"func":   "handleFollow",
		"lineID": accountID,
	}).Info("New Line Follow")
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
