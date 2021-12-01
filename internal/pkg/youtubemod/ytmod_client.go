package youtubemod

import (
	"encoding/xml"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
)

type YoutubeClient struct {
	Client      *Client
	CallbackUrl string
}

var YC = YoutubeClient{}

var (
	channelBaseUrl = "https://www.youtube.com/xml/feeds/videos.xml?channel_id="
	pubSubBaseUrl  = "https://pubsubhubbub.appspot.com/subscribe"
)

func (yc *YoutubeClient) Init(host string, psPort int) {
	YC.Client = NewClient(pubSubBaseUrl, host, psPort, "test app")
}

func (yc *YoutubeClient) CreatePubSubByChannelId(channelId string) {
	yc.Client.Subscribe(channelBaseUrl+channelId, FeedHandler)
}

type Feed struct {
	Status  string  `xml:"status>http"`
	Xmlns   string  `xml:"-xmlns"`
	Yt      string  `xml:"-yt"`
	Title   string  `xml:"title"`
	Updated string  `xml:"updated"`
	Entries []Entry `xml:"entry" validate:"dive"`
}

type Entry struct {
	VideoID   string   `xml:"videoId" validate:"required"`
	ChannelId string   `xml:"channelId" validate:"required"`
	Published string   `xml:"published"`
	Updated   string   `xml:"updated"`
	Title     string   `xml:"title" validate:"required"`
	EntryLink []Link   `xml:"link"`
	Author    []Author `xml:"author"`
}

type Author struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

func FeedHandler(contentType string, body []byte) {
	var feed Feed

	fmt.Println("Get PubSub feed: ", contentType)

	xmlError := xml.Unmarshal(body, &feed)
	if xmlError != nil {
		log.Printf("XML Parse Error %v", xmlError)
	} else {
		fmt.Printf("feed: %+v\n", feed)
	}

	//Validate the feed
	validate := validator.New()
	validErr := validate.Struct(feed)
	if validErr != nil {
		fmt.Println("Invalid Feed from pubSub:", validErr)
		return
	}
	Tracker.AddUpcoming(feed)
}
