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

	testUser = model.User{
		Line:            testLine,
		LineAccessToken: testAccessToken,
	}

	TestChannelIds = []string{
		"UC1DCedRgGHBdm81E1llLhOQ",
		"UC-hM6YJuNYVAmUWxeIr9FeA",
		"UC1opHUrw8rvnsadT-iGp7Cg",
		"UCCzUftO8KOVkV4wQG1vkUvg",
		"UCl_gCybOJRIgOXw6Qb4qJzQ",
		"UCiEm9noegBIb-AzjqpxKffA", //羅傑
		"UCqm3BQLlJfvkTsX_hvm0UmA", //WTM
		"UCMwGHR0BTZuLsmjY_NT5Pwg", //Ina
		"UChgTyjG-pdNvxxhdsXfHQ5Q", //Pavolia
		"UCD8HOxPs4Xvsm8H0ZxXGiBw", //Mel
		"UC_vMYWcDjmfdpH6r4TTn1MQ", //Iroha
		"UCvInZx9h3jC2JzsIzoOebWg", //Flare
		"UC4G-xDOf5U9luBcfpyaqF3Q", //My Channel
		"UCZlDXzGoo7d44bwdNObFacg", //Katana
		"UCK9V2B22uJYu3N7eR_BT9QA", //polka
	}

	TestVideoIds = []string{
		"6hZ-kf1aQ1M",
		"omgSWqwVTjY",
		"IwlECRC8c0E",
		"D2tjfs_Dn7A",
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
	for i, channelId := range TestChannelIds {
		youtubemod.CreatePubSubByChannelId(channelId)
		model.DB.CreateYtSubscription(&model.YtSubscription{
			Line:            testLine,
			LineAccessToken: testAccessToken,
			ChannelId:       channelId,
		})
		model.DB.CreatePubSubSubscription(&model.PubSubSubscription{
			Topic:      "https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelId,
			CallbackId: i,
		})
	}
	for _, id := range TestVideoIds {
		model.DB.CreateYtVideo(&model.YtVideo{
			VideoId: id,
		})
	}
	// Restore PubSubClient from DB
	var oldPubSubs []model.PubSubSubscription
	oldPubSubs = model.DB.QueryAllPubSub()
	for _, old := range oldPubSubs {
		youtubemod.PubSub.RestoreSubscribe(old.Topic, old.CallbackId, youtubemod.FeedHandler)
	}

	//step 2.2.1: init YouTube Tracker
	youtubemod.Tracker.Init()
	go youtubemod.Tracker.StartTrack()

	http.HandleFunc("/youtube/pubsub/", youtubemod.PubSub.HandlePubSubCallback)

	//step 2.2.2: create test PubSub

	//step 3: start up our webhook server
	fmt.Println("Starting the webserver listen on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
