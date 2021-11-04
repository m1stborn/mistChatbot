package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/pubsubhubbub/gohubbub"
)

func main() {

	var (
		host   = "popular-dolphin-4.loca.lt"
		psPort = 2323
	)

	ID := []string{"UC1DCedRgGHBdm81E1llLhOQ", "UC-hM6YJuNYVAmUWxeIr9FeA", "UC1opHUrw8rvnsadT-iGp7Cg", "UCCzUftO8KOVkV4wQG1vkUvg"}

	psClient := gohubbub.NewClient("https://pubsubhubbub.appspot.com/subscribe", host, psPort, "test app")
	for _, id := range ID {
		topic := fmt.Sprintf("https://www.youtube.com/xml/feeds/videos.xml?channel_id=%s", id)
		psClient.Subscribe(topic, FeedHandler3)
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

func FeedHandler3(contentType string, body []byte) {
	var feed map[string]interface{}
	fmt.Println("Get a feed: ", contentType)
	xmlError := xml.Unmarshal(body, &feed)

	if xmlError != nil {
		log.Printf("XML Parse Error %v", xmlError)

	} else {
		fmt.Println(feed)
	}
}
