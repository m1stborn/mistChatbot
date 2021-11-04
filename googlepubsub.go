package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/m1stborn/mistChatbot/internal/pkg/youtubemod"
)

func main() {

	var (
		host   = "rude-earwig-3.loca.lt"
		psPort = 1919
	)

	ID := []string{
		"UC1DCedRgGHBdm81E1llLhOQ",
		"UC-hM6YJuNYVAmUWxeIr9FeA",
		"UC1opHUrw8rvnsadT-iGp7Cg",
		"UCCzUftO8KOVkV4wQG1vkUvg",
		"UCl_gCybOJRIgOXw6Qb4qJzQ",
	}
	//ID := []string{"UC1opHUrw8rvnsadT-iGp7Cg", "UC4G-xDOf5U9luBcfpyaqF3Q"}

	psClient := youtubemod.NewClient("https://pubsubhubbub.appspot.com/subscribe", host, psPort, "test app")
	for _, id := range ID {
		topic := fmt.Sprintf("https://www.youtube.com/xml/feeds/videos.xml?channel_id=%s", id)
		psClient.Subscribe(topic, FeedHandler)
	}

	go psClient.StartServer()

	time.Sleep(time.Second * 5)
	log.Println("Press Enter for graceful shutdown...")

	var input string
	fmt.Scanln(&input)

	for _, id := range ID {
		topic := fmt.Sprintf("https://www.youtube.com/xml/feeds/videos.xml?channel_id=%s", id)
		psClient.Unsubscribe(topic)
	}

	time.Sleep(time.Second * 5)
}

type Feed struct {
	Status  string  `xml:"status>http"`
	Xmlns   string  `xml:"-xmlns"`
	Yt      string  `xml:"-yt"`
	Title   string  `xml:"title"`
	Updated string  `xml:"updated"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	VideoID   string   `xml:"videoId"`
	ChannelId string   `xml:"channelId"`
	Published string   `xml:"published"`
	Updated   string   `xml:"updated"`
	Title     string   `xml:"title"`
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

type Map map[string]interface{}

func FeedHandler(contentType string, body []byte) {
	var feed Feed

	fmt.Println("Get a feed: ", contentType)

	xmlError := xml.Unmarshal(body, &feed)

	if xmlError != nil {
		log.Printf("XML Parse Error %v", xmlError)

	} else {
		fmt.Printf("feed: %+v\n", feed)
	}

	//fmt.Println("Get a feed: ", contentType)
	//mv, xmlErr := mxj.NewMapXml(body)
	//if xmlErr != nil {
	//	log.Printf("XML Parse Error %v", xmlErr)
	//} else {
	//	fmt.Println(mv)
	//}

}
