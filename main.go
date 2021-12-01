package main

import (
	"fmt"
	"github.com/m1stborn/mistChatbot/internal/pkg/youtubemod"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/m1stborn/mistChatbot/internal/pkg/line"
	"github.com/m1stborn/mistChatbot/internal/pkg/model"
	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"

	_ "github.com/joho/godotenv/autoload"
)

var err error

var (
	port   = os.Getenv("PORT")
	dbUri  = os.Getenv("DB_URI")
	psHost = os.Getenv("CALLBACK_URL_BASE")

	testStreamer    = []string{"muse_tw", "lck", "dogdog", "lolpacifictw", "m989876525", "qq7925168", "never_loses"}
	testLine        = os.Getenv("TEST_LINE_USER")
	testAccessToken = os.Getenv("TEST_LINE_NOTIFY_ACCESSTOKEN")

	testUser = model.User{
		Line:            testLine,
		LineAccessToken: testAccessToken,
	}
)

func main() {

	//step 1: init DB
	model.DB.Init(dbUri)

	model.DB.CreateUser(&testUser)

	//step 1.1: Create http router for line webhook
	http.HandleFunc("/line", line.Handler)
	http.HandleFunc("/line/notify/auth", line.HandelNotifyAuth)
	http.HandleFunc("/line/notify/callback", line.HandleNotifyCallback)

	//step 1.2: Create http router for twitch webhook
	http.HandleFunc("/twitch/channelFollow", twitchmod.EventSubFollow)
	http.HandleFunc("/twitch/streamOnline", twitchmod.EventSubStreamOnline)

	//step 2.1: init Twitch Client
	//twitch := nweTwitch()

	//step 2.1.1: delete old subscription during develop and testing
	subIds := twitchmod.GetSubscriptions()
	twitchmod.DeleteSubscriptions(subIds)

	//step 2.1.2: Create Event subscriptions
	err = twitchmod.CreateChannelFollowSubscription("twitch", "/twitch/channelFollow")
	err = twitchmod.CreateStreamOnlineSubscription("twitch", "/twitch/streamOnline")

	for _, streamer := range testStreamer {
		err = twitchmod.CreateStreamOnlineSubscription(streamer, "/twitch/streamOnline")
		model.DB.CreateSubscription(&model.Subscription{
			Line:            testLine,
			LineAccessToken: testAccessToken,
			TwitchLoginName: streamer,
		})
	}

	if err != nil {
		fmt.Println(err)
	}

	portInt, portErr := strconv.Atoi(port)
	if portErr != nil {
		fmt.Println(err)
	}

	//step 2.2.1: init YouTube Client
	youtubemod.Tracker.Init()
	youtubemod.YC.Init(psHost, portInt)

	//step 2.2.2: start up PubSub client and Tracker
	go youtubemod.YC.Client.StartServer()
	go youtubemod.Tracker.StartTrack()

	//step 2.2.2: create test PubSub
	for _, channelId := range TestChannelId {
		youtubemod.YC.CreatePubSubByChannelId(channelId)
		model.DB.CreateYtSubscription(&model.YtSubscription{
			Line:            testLine,
			LineAccessToken: testAccessToken,
			ChannelId:       channelId,
		})
	}

	//step 3: start up our webhook server
	fmt.Println("Starting the webserver listen on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
