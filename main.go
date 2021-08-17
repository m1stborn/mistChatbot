package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/m1stborn/mistChatbot/internal/pkg/line"
	"github.com/m1stborn/mistChatbot/internal/pkg/model"
	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nicklaw5/helix"
)

var (
	//lineSecretToken = os.Getenv("LINE_CHANNEL_SECRET")
	//lineAccessToken = os.Getenv("LINE_CHANNEL_ACCESSTOKEN")

	twitchClientID    = os.Getenv("TWITCH_CLIENT_ID")
	twitchAccessToken = os.Getenv("TWITCH_ACCESSTOKEN")
	secretWord        = "s3cre7w0rd"

	callbackUrl = os.Getenv("CALLBACK_URL")
	port        = ":" + os.Getenv("PORT")

	//lineClient *linebot.Client
	//err        error

	dbUri = os.Getenv("DB_URI")

	testStreamer    = []string{"muse_tw", "lck", "dogdog", "lolpacifictw", "m989876525", "qq7925168", "never_loses"}
	testLine        = os.Getenv("TEST_LINE_USER")
	testAccessToken = os.Getenv("TEST_LINE_NOTIFY_ACCESSTOKEN")

	testUser = model.User{
		Line:            testLine,
		LineAccessToken: testAccessToken,
	}
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
	//lineClient, err = linebot.New(lineSecretToken, lineAccessToken)
	//if err != nil {
	//	log.Println(err.Error())
	//}

	//step 1: init DB
	model.DB.Init(dbUri)

	model.DB.CreateUser(&testUser)

	//step 1.1: Create http router for line webhook
	http.HandleFunc("/line", line.Handler)
	http.HandleFunc("/line/notify/auth", line.HandelNotifyAuth)
	http.HandleFunc("/line/notify/callback", line.HandleNotifyCallback)

	//step 2: init Twitch Client
	twitch := nweTwitch()

	//step 2.0: delete old subscription during develop and testing
	subIds := twitch.GetSubscriptions()
	twitch.DeleteSubscriptions(subIds)

	//step 2.1: Create Event subscriptions
	twitch.CreateChannelFollowSubscription("twitch", "/callback/channelFollow")
	twitch.CreateStreamOnlineSubscription("twitch", "/callback/streamOnline")

	for _, streamer := range testStreamer {
		twitch.CreateStreamOnlineSubscription(streamer, "/callback/streamOnline")
		model.DB.CreateSubscription(&model.Subscription{
			LineUser:        testLine,
			TwitchLoginName: streamer,
		})
	}

	//step 2.2: Create http router for twitch webhook
	http.HandleFunc("/callback/channelFollow", twitchmod.EventSubFollow)
	http.HandleFunc("/callback/streamOnline", twitchmod.EventSubStreamOnline)

	//step 3: start up our webhook server
	fmt.Println("Starting the webserver listen on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
