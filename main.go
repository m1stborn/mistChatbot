package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nicklaw5/helix"
)

const (
	lineSecretToken = "4fabd1de4303c0dfc00999e0200a9438"
	lineAccessToken = "2WjKzdmNn/lpLHSa0Yv+G50sBrV7gvTg7hqbqZS+wpfVJg2fqYmwFWWxtBkBMjl2KZtJuAhCXXds7lqlCcQyVhVozxloEh3UTOwnWp5km735r6hT2f2zMDG7Av7mXmcJq/HqJABeagd5f9IQRyydQwdB04t89/1O/w1cDnyilFU="

	twitchClientID    = "uq2sfwwume0bcdxsrfr7yvhlm8omkg"
	twitchAccessToken = "iatwfdt5hws5ggusoe7zf09odx0hsk"
	secretWord        = "s3cre7w0rd"

	callbackUrl = "https://salty-ocean-83656.herokuapp.com"
)

var (
	lineClient *linebot.Client
	err        error

	port = ":" + os.Getenv("PORT")
)

func nweTwitch() twitchmod.TwitchClient {
	cl, clErr := helix.NewClient(&helix.Options{
		ClientID:       twitchClientID,
		AppAccessToken: twitchAccessToken,
	})
	//TODO integrate logging and handle error
	if clErr != nil {
		fmt.Println()
	}
	twitch := twitchmod.TwitchClient{
		Client:      cl,
		CallbackUrl: callbackUrl,
		SecretWord:  secretWord,
	}
	return twitch
}

func main() {
	//step 1: init Line Client
	lineClient, err = linebot.New(lineSecretToken, lineAccessToken)
	//TODO integrate logging and handle error
	if err != nil {
		log.Println(err.Error())
	}

	//step 1.1: Create http router for line webhook
	http.HandleFunc("/callback", callbackHandler)

	//step 2: init Twitch Client
	twitch := nweTwitch()

	//step 2.0: delete old subscription during develop and testing
	subIds := twitch.GetSubscriptions()
	fmt.Println("event subID to delete:", subIds)

	twitch.DeleteSubscriptions(subIds)

	//step 2.1: Create Event subscriptions
	twitch.CreateChannelFollowSubscription("twitch", "/callback/channelFollow")
	twitch.CreateStreamOnlineSubscription("twitch", "/callback/streamOnline")

	//step 2.2: Create http router for twitch webhook
	http.HandleFunc("/callback/channelFollow", twitchmod.EventSubFollow)
	http.HandleFunc("/callback/streamOnline", twitchmod.EventSubStreamOnline)

	//step 3: start up our webhook server
	fmt.Println("Starting the webserver listen on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, cbErr := lineClient.ParseRequest(r)

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
				if _, cbErr = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}
