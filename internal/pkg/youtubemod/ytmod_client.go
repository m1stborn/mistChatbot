package youtubemod

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var (
	port = os.Getenv("PORT")
	host = os.Getenv("CALLBACK_URL_BASE")

	channelBaseUrl = "https://www.youtube.com/xml/feeds/videos.xml?channel_id="
)

var PubSub = &PubSubClient{}

func init() {
	portInt, portErr := strconv.Atoi(port)
	if portErr != nil {
		fmt.Println(portErr)
	}

	PubSub = NewPubSubClient(host, portInt, "test app")
	PubSub.StartClient()
}

func CreatePubSubByChannelId(channelId string) {
	PubSub.Subscribe(channelBaseUrl+channelId, FeedHandler)
}

func UnsubscribePubSubByChannelId(channelId string) {
	PubSub.Unsubscribe(channelBaseUrl + channelId)
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
