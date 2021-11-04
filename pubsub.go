package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/dpup/gohubbub"
)

func main() {

	var (
		host   = "tender-snake-59.loca.lt"
		psPort = 2121
	)

	ID := []string{"UC1DCedRgGHBdm81E1llLhOQ", "UC-hM6YJuNYVAmUWxeIr9FeA", "UC1opHUrw8rvnsadT-iGp7Cg", "UCCzUftO8KOVkV4wQG1vkUvg"}

	psClient := gohubbub.NewClient(fmt.Sprintf("%s:%d", host, psPort), "test app")
	for _, id := range ID {
		topic := fmt.Sprintf("https://www.youtube.com/xml/feeds/videos.xml?channel_id=%s", id)
		psErr := psClient.DiscoverAndSubscribe(topic, HandleFeed)
		if psErr != nil {
			fmt.Println(psErr)
		}
	}

	//psErr := psClient.DiscoverAndSubscribe("https://www.youtube.com/xml/feeds/videos.xml?channel_id=UCc5bWVuWyL74ngKNjm49hdw", HandlerFeed)
	//if psErr != nil {
	//	fmt.Println(psErr)
	//}

	//go psClient.StartAndServe("", psPort)
	psClient.StartAndServe("", psPort)
	//msg := "{\"assets\" : {\"old\" : 123}}"
	//var payload map[string]interface{}
	//var payload2 json.RawMessage
	//jsonErr := json.Unmarshal([]byte(msg), &payload2)
	//if jsonErr != nil {
	//	fmt.Println(jsonErr)
	//}
	//fmt.Println(payload)
}

func HandleFeed(contentType string, body []byte) {
	var feed map[string]interface{}
	fmt.Println("Get a feed: ", contentType)
	xmlError := xml.Unmarshal(body, &feed)

	if xmlError != nil {
		log.Printf("XML Parse Error %v", xmlError)

	} else {
		fmt.Println(feed)
	}
}
