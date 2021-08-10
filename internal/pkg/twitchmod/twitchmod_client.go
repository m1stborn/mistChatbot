package twitchmod

import (
	"fmt"

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

	//TODO integrate logging and handle error
	if err != nil {
		//fmt.Println("GetSubscriptions err:", err)
		logger.WithField("func", "GetSubscriptions").Error(err.Error())
	}
	//fmt.Printf("GetSubscriptions resp:%+v\n", resp)
	logger.WithField("func", "GetSubscriptions").Infof("%+v\n", resp)

	if resp != nil {
		for _, data := range resp.Data.EventSubSubscriptions {
			fmt.Println(data)
			idList = append(idList, data.ID)
		}
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

	//TODO integrate logging and handle error
	if err != nil {
		logger.WithField("func", "CreateChannelFollowSubscription").Error(err.Error())
	}
	//fmt.Printf("CreateChannelFollowSubscription resp:%+v\n", resp)
	logger.WithField("func", "CreateChannelFollowSubscription").Infof("%+v\n", resp)
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

	//TODO integrate logging and handle error
	if err != nil {
		//fmt.Println("CreateStreamOnlineSubscription err:", err)
		logger.WithField("func", "CreateEventSubSubscription").Error(err.Error())
	}
	//fmt.Printf("CreateStreamOnlineSubscription resp:%+v\n", resp)
	logger.WithField("func", "CreateEventSubSubscription").Infof("%+v\n", resp)
}

func (c *TwitchClient) DeleteSubscriptions(idList []string) {
	for _, id := range idList {
		deleteResp, deleteErr := c.Client.RemoveEventSubSubscription(id)

		//TODO integrate logging and handle error
		if deleteErr != nil {
			logger.WithField("func", "DeleteSubscriptions").Error(deleteErr.Error())
		}
		fmt.Printf("DeleteSubscriptions:%+v\n", deleteResp)
	}
}

func (c *TwitchClient) GetUsersID(usernameList []string) (idList []string) {
	userResp, userErr := c.Client.GetUsers(&helix.UsersParams{
		//IDs:    []string{"twitch user id"},
		//Logins: []string{"twitch user name"},
		Logins: usernameList,
	})

	//TODO integrate logging and handle error
	if userErr != nil {
		//fmt.Println("GetUsersID error:", userErr)
		logger.WithField("func", "GetUsersID").Error(userErr.Error())
	}
	//fmt.Printf("GetUsersID resp:%+v\n", userResp)
	logger.WithField("func", "GetUsersID").Infof("%+v\n", userResp)

	if userResp != nil {
		for _, user := range userResp.Data.Users {
			idList = append(idList, user.ID)
		}
	}
	return idList
}
