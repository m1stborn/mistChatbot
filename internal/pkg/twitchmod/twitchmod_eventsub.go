package twitchmod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nicklaw5/helix"
)

type EventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func EventSubFollow(w http.ResponseWriter, r *http.Request) {

	//TODO integrate logging
	fmt.Println("Receive a http request:", r.Method)

	body, errF := ioutil.ReadAll(r.Body)
	if errF != nil {
		log.Println(errF)
		return
	}
	defer r.Body.Close()
	// verify that the notification came from twitch using the secret.
	if !helix.VerifyEventSubNotification("s3cre7w0rd", r.Header, string(body)) {
		log.Println("no valid signature on subscription")
		return
	} else {
		log.Println("verified signature for subscription")
	}
	var vals EventSubNotification
	errF = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if errF != nil {
		log.Println(errF)
		return
	}
	// if there's a challenge in the request, respond with only the challenge to verify your eventsub.
	if vals.Challenge != "" {
		w.Write([]byte(vals.Challenge))
		return
	}
	var followEvent helix.EventSubChannelFollowEvent
	errF = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&followEvent)

	log.Printf("got follow webhook: %s follows %s\n", followEvent.UserName, followEvent.BroadcasterUserName)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func EventSubStreamOnline(w http.ResponseWriter, r *http.Request) {

	//TODO integrate logging
	fmt.Println("Receive a http request:", r.Method)

	body, errF := ioutil.ReadAll(r.Body)
	if errF != nil {
		log.Println(errF)
		return
	}
	defer r.Body.Close()
	// verify that the notification came from twitch using the secret.
	if !helix.VerifyEventSubNotification("s3cre7w0rd", r.Header, string(body)) {
		log.Println("no valid signature on subscription")
		return
	} else {
		log.Println("verified signature for subscription")
	}
	var vals EventSubNotification
	errF = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if errF != nil {
		log.Println(errF)
		return
	}
	// if there's a challenge in the request, respond with only the challenge to verify your eventsub.
	if vals.Challenge != "" {
		w.Write([]byte(vals.Challenge))
		return
	}

	var streamOnlineEvent helix.EventSubStreamOnlineEvent
	errF = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamOnlineEvent)

	log.Printf("got stream online webhook: %s start stream\n", streamOnlineEvent.BroadcasterUserName)
	w.WriteHeader(200)
	w.Write([]byte("ok"))

	//stream online notification latency about 1 min
}
