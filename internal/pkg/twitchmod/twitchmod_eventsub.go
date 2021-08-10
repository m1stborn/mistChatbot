package twitchmod

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

type EventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func EventSubFollow(w http.ResponseWriter, r *http.Request) {

	logger.WithFields(log.Fields{
		"func":   "EventSubFollow",
		"method": r.Method,
	})

	body, errF := ioutil.ReadAll(r.Body)
	if errF != nil {
		logger.WithField("func", "EventSubFollow").Error(errF.Error())
		return
	}
	defer r.Body.Close()
	// verify that the notification came from twitch using the secret.
	if !helix.VerifyEventSubNotification("s3cre7w0rd", r.Header, string(body)) {
		logger.WithField("func", "EventSubFollow").Info("no valid signature on subscription")
		return
	} else {
		logger.WithField("func", "EventSubFollow").Info("verified signature for subscription")
	}
	var vals EventSubNotification
	errF = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if errF != nil {
		logger.WithField("func", "EventSubFollow").Error(errF.Error())
		return
	}
	// if there's a challenge in the request, respond with only the challenge to verify your eventSub.
	if vals.Challenge != "" {
		w.Write([]byte(vals.Challenge))

		return
	}
	var followEvent helix.EventSubChannelFollowEvent
	errF = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&followEvent)

	logger.WithFields(log.Fields{
		"func": "EventSubFollow",
		"type": "Notify",
	}).Infof("%s follows %s!\n", followEvent.UserName, followEvent.BroadcasterUserName)

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func EventSubStreamOnline(w http.ResponseWriter, r *http.Request) {

	logger.WithFields(log.Fields{
		"func":   "EventSubStreamOnline",
		"method": r.Method,
	})

	body, errF := ioutil.ReadAll(r.Body)
	if errF != nil {
		logger.WithField("func", "EventSubStreamOnline").Error(errF.Error())
		return
	}
	defer r.Body.Close()
	// verify that the notification came from twitch using the secret.
	if !helix.VerifyEventSubNotification("s3cre7w0rd", r.Header, string(body)) {
		logger.WithField("func", "EventSubStreamOnline").Info("no valid signature on subscription")
		return
	} else {
		logger.WithField("func", "EventSubStreamOnline").Info("verified signature for subscription")
	}
	var vals EventSubNotification
	errF = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
	if errF != nil {
		logger.WithField("func", "EventSubStreamOnline").Error(errF.Error())
		return
	}
	// if there's a challenge in the request, respond with only the challenge to verify your eventsub.
	if vals.Challenge != "" {
		w.Write([]byte(vals.Challenge))
		return
	}

	var streamOnlineEvent helix.EventSubStreamOnlineEvent
	errF = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&streamOnlineEvent)

	logger.WithFields(log.Fields{
		"func": "EventSubFollow",
		"type": "Notify",
	}).Infof("%s start streaming!\n", streamOnlineEvent.BroadcasterUserName)

	w.WriteHeader(200)
	w.Write([]byte("ok"))

	//stream online notification latency about 1 min
}
