package line

import (
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	lineSecretToken = "4fabd1de4303c0dfc00999e0200a9438"
	lineAccessToken = "2WjKzdmNn/lpLHSa0Yv+G50sBrV7gvTg7hqbqZS+wpfVJg2fqYmwFWWxtBkBMjl2KZtJuAhCXXds7lqlCcQyVhVozxloEh3UTOwnWp5km735r6hT2f2zMDG7Av7mXmcJq/HqJABeagd5f9IQRyydQwdB04t89/1O/w1cDnyilFU="
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
