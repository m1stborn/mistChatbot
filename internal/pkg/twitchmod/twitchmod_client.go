package twitchmod

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nicklaw5/helix"
)

type TwitchClient struct {
	Client            *helix.Client
	CallbackUrl       string
	SecretWord        string
	twitchClientID    string
	twitchAccessToken string
}

//var Twitch = TwitchClient{}
//
//func init() {
//	var err error
//	Twitch.client, err = helix.NewClient(&helix.Options{
//		ClientID:       twitchClientID,
//		AppAccessToken: twitchAccessToken,
//	})
//	if err != nil {
//		//handle error
//	}
//	Twitch.callbackUrl = "https://red-panda-59.loca.lt" + "/callback"
//}

func (c *TwitchClient) GetSubscriptions() (idList []string) {
	resp, err := c.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		Status: helix.EventSubStatusEnabled, // This is optional
	})

	if err != nil {
		logger.WithField("func", "GetSubscriptions").Error(err.Error())
	}
	if resp != nil {
		//TODO: logging whole response?
		//logger.WithFields(log.Fields{
		//	"func": "GetSubscriptions",
		//}).Infof(fmt.Sprintf("resp.data: %+v", resp.Data))

		//var respList []struct {
		//	ID string
		//	TYPE string
		//}
		for _, data := range resp.Data.EventSubSubscriptions {
			idList = append(idList, data.ID)
			//respList = append(respList, struct {
			//	ID string
			//	TYPE string
			//}{data.ID,data.Type})
		}
		logger.WithFields(log.Fields{
			"func": "GetSubscriptions",
		}).Infof(fmt.Sprintf("Subscriptions ID List: %+v", idList))
	}

	return idList //return current eventSubscriptions id
}

func (c *TwitchClient) CreateChannelFollowSubscription(broadcasterName string, route string) {
	idList := c.GetUsersID([]string{broadcasterName})
	id := idList[0]

	resp, err := c.Client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeChannelFollow,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: c.CallbackUrl + route,
			Secret:   c.SecretWord,
		},
	})

	if err != nil {
		logger.WithField("func", "CreateChannelFollowSubscription").Error(err.Error())
	}
	if resp != nil {
		//TODO: logging whole response?
		logger.WithFields(log.Fields{
			"func": "CreateChannelFollowSubscription",
		}).Infof(fmt.Sprintf("broadcaster ID: %+v", id))
	}
}

func (c *TwitchClient) CreateStreamOnlineSubscription(broadcasterName string, route string) {
	idList := c.GetUsersID([]string{broadcasterName})
	id := idList[0]

	resp, err := c.Client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeStreamOnline,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: c.CallbackUrl + route,
			Secret:   c.SecretWord,
		},
	})

	if err != nil {
		logger.WithField("func", "CreateStreamOnlineSubscription").Error(err.Error())
	}
	if resp != nil {
		//TODO: logging whole response?
		logger.WithFields(log.Fields{
			"func": "CreateStreamOnlineSubscription",
		}).Infof(fmt.Sprintf("broadcaster ID: %+v", id))

	}
}

func (c *TwitchClient) DeleteSubscriptions(idList []string) {
	for _, id := range idList {
		deleteResp, deleteErr := c.Client.RemoveEventSubSubscription(id)

		if deleteErr != nil {
			logger.WithField("func", "DeleteSubscriptions").Error(deleteErr.Error())
		}
		if deleteResp != nil {
			logger.WithFields(log.Fields{
				"func": "DeleteSubscriptions",
			}).Infof(fmt.Sprintf("deleteID: %+v", id))
		}
	}
}

func (c *TwitchClient) GetUsersID(usernameList []string) (idList []string) {
	userResp, userErr := c.Client.GetUsers(&helix.UsersParams{
		//example usage
		//IDs:    []string{"twitch user id"},
		//Logins: []string{"twitch user name"},
		Logins: usernameList,
	})

	if userErr != nil {
		logger.WithField("func", "GetUsersID").Error(userErr.Error())
	}
	if userResp != nil {
		//TODO: logging whole response?
		//logger.WithFields(log.Fields{
		//	"func": "GetUsersID",
		//}).Infof(fmt.Sprintf("resp.data: %+v", userResp.Data))

		for _, user := range userResp.Data.Users {
			idList = append(idList, user.ID)
		}

		logger.WithFields(log.Fields{
			"func": "GetUsersID",
		}).Infof(fmt.Sprintf("User ID List: %+v", idList))
	}

	return idList
}