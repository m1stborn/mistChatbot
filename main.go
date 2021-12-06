package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/m1stborn/mistChatbot/internal/pkg/line"
	"github.com/m1stborn/mistChatbot/internal/pkg/model"
	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"
	"github.com/m1stborn/mistChatbot/internal/pkg/youtubemod"

	_ "github.com/joho/godotenv/autoload"
)

var err error

var (
	port  = os.Getenv("PORT")
	dbUri = os.Getenv("DB_URI")
)

var (
	testStreamer    = []string{"muse_tw", "lck", "dogdog", "lolpacifictw", "m989876525", "qq7925168", "never_loses"}
	testLine        = os.Getenv("TEST_LINE_USER")
	testAccessToken = os.Getenv("TEST_LINE_NOTIFY_ACCESSTOKEN")
)

func main() {

	//step 1: init DB
	model.DB.Init(dbUri)

	//step 1.1: Create http router for line webhook
	http.HandleFunc("/line", line.Handler)
	http.HandleFunc("/line/notify/auth", line.HandelNotifyAuth)
	http.HandleFunc("/line/notify/callback", line.HandleNotifyCallback)

	//step 1.2: Create http router for twitch webhook
	http.HandleFunc("/twitch/channelFollow", twitchmod.EventSubFollow)
	http.HandleFunc("/twitch/streamOnline", twitchmod.EventSubStreamOnline)

	//step 1.2: Create http router for YouTube PubSubHubBub webhook
	http.HandleFunc("/youtube/pubsub/", youtubemod.PubSub.HandlePubSubCallback)

	//step 2.1: init Twitch Client
	//twitch := nweTwitch()

	//step 2.2: delete old subscription during develop and testing
	//subIds := twitchmod.GetSubscriptions()
	//twitchmod.DeleteSubscriptions(subIds)

	//step 2.3: create Event subscriptions
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

	//step 3.1: restore PubSubClient from DB
	var oldPubSubs []model.PubSubSubscription
	oldPubSubs = model.DB.QueryAllPubSub()
	for _, old := range oldPubSubs {
		youtubemod.PubSub.RestoreSubscribe(old.Topic, old.CallbackId, youtubemod.FeedHandler)
	}

	//step 3.2: init YouTube Tracker
	youtubemod.Tracker.Init()
	go youtubemod.Tracker.StartTrack()

	//step 4: start up server
	fmt.Println("Starting the webserver listen on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
