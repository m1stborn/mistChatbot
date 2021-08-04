package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	secretToken = "4fabd1de4303c0dfc00999e0200a9438"
	accessToken = "2WjKzdmNn/lpLHSa0Yv+G50sBrV7gvTg7hqbqZS+wpfVJg2fqYmwFWWxtBkBMjl2KZtJuAhCXXds7lqlCcQyVhVozxloEh3UTOwnWp5km735r6hT2f2zMDG7Av7mXmcJq/HqJABeagd5f9IQRyydQwdB04t89/1O/w1cDnyilFU="
)

var (
	client *linebot.Client
	err    error
)

func main() {
	fmt.Println("Hello, this is chat bot!")

	client, err = linebot.New(secretToken, accessToken)

	if err != nil {
		log.Println(err.Error())
	}

	http.HandleFunc("/callback", callbackHandler)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := client.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
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
				if _, err = client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}
