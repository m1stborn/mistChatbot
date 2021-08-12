package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/m1stborn/mistChatbot/internal/pkg/line"
	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"

	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/nicklaw5/helix"
)

var (
	lineSecretToken = os.Getenv("LINE_CHANNEL_SECRET")
	lineAccessToken = os.Getenv("LINE_CHANNEL_ACCESSTOKEN")

	twitchClientID    = os.Getenv("TWITCH_CLIENT_ID")
	twitchAccessToken = os.Getenv("TWITCH_ACCESSTOKEN")
	secretWord        = "s3cre7w0rd"

	callbackUrl = os.Getenv("CALLBACK_URL")
	port        = ":" + os.Getenv("PORT")
)

var (
	lineClient *linebot.Client
	err        error
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
	http.HandleFunc("/line", line.Handler)

	//step 2: init Twitch Client
	twitch := nweTwitch()

	//step 2.0: delete old subscription during develop and testing
	subIds := twitch.GetSubscriptions()
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
