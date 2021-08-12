package line

import (
	"fmt"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
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
		fmt.Println(err)
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

		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, cbErr = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
