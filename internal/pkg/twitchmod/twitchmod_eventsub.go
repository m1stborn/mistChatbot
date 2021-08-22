package twitchmod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/m1stborn/mistChatbot/internal/pkg/model"

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
	if !helix.VerifyEventSubNotification(secretWord, r.Header, string(body)) {
		logger.WithField("func", "EventSubFollow").Info("no valid signature on subscription")
		return
	} else {
		//logger.WithField("func", "EventSubFollow").Info("verified signature for subscription")
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
	}).Infof("%s follows %s!", followEvent.UserName, followEvent.BroadcasterUserName)

	w.WriteHeader(200)
	w.Write([]byte("ok"))

	//SendLineNotify(testAccessToken,
	//	fmt.Sprintf("%s follows %s!",
	//		followEvent.UserName,
	//		followEvent.BroadcasterUserName))

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
	if !helix.VerifyEventSubNotification(secretWord, r.Header, string(body)) {
		logger.WithField("func", "EventSubStreamOnline").Info("no valid signature on subscription")
		return
	} else {
		//logger.WithField("func", "EventSubStreamOnline").Info("verified signature for subscription")
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
	}).Infof("%s start streaming!", streamOnlineEvent.BroadcasterUserName)

	w.WriteHeader(200)
	w.Write([]byte("ok"))

	//stream online notification latency about 1 min

	accessTokens := model.DB.QuerySubByTwitchLoginName(streamOnlineEvent.BroadcasterUserLogin)

	for _, token := range accessTokens {
		SendLineNotify(token,
			fmt.Sprintf("%s start streaming!\n https://www.twitch.tv/%s",
				streamOnlineEvent.BroadcasterUserName,
				streamOnlineEvent.BroadcasterUserLogin))
	}

}

const notifyAPIHost string = "https://notify-api.line.me"

func SendLineNotify(accessToken string, message string) {
	uri := "/api/notify"
	queryStr := url.Values{}
	queryStr.Add("message", message)
	encodeQueryStr := queryStr.Encode()
	pr, httpErr := http.NewRequest("POST", notifyAPIHost+uri, bytes.NewBufferString(encodeQueryStr))
	pr.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	pr.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	r, httpErr := client.Do(pr)
	if httpErr != nil {
		log.WithError(httpErr).Error("Notify Request Failed")
		return
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(r.Body)
		logger.WithFields(log.Fields{
			"status":   r.Status,
			"response": string(data),
		}).Error("LINE Notify Failed")
	}
}
