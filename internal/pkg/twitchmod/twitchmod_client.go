package twitchmod

import (
	"fmt"

	"github.com/nicklaw5/helix"
)

type TwitchClient struct {
	client            *helix.Client
	callbackUrl       string
	secretWord        string
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
	resp, err := c.client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		Status: helix.EventSubStatusEnabled, // This is optional
	})

	//TODO integrate logging and handle error
	if err != nil {
	}
	if resp != nil {
		for _, data := range resp.Data.EventSubSubscriptions {
			fmt.Println(data)
			idList = append(idList, data.ID)
		}
	}

	return idList //return current eventSubscriptions id
}

func (c *TwitchClient) CreateChannelFollowSubscription(broadcasterName string) {
	idList := c.GetUsersID([]string{broadcasterName})
	id := idList[0]

	resp, err := c.client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeChannelFollow,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: c.callbackUrl,
			Secret:   c.secretWord,
		},
	})

	//TODO integrate logging and handle error
	if err != nil {
		fmt.Println("eventSubErr:", err)
	}
	fmt.Printf("%+v\n", resp)
}

func (c *TwitchClient) CreateStreamOnlineSubscription(broadcasterName string) {
	idList := c.GetUsersID([]string{broadcasterName})
	id := idList[0]

	resp, err := c.client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeStreamOnline,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: c.callbackUrl,
			Secret:   c.secretWord,
		},
	})

	//TODO integrate logging and handle error
	if err != nil {
		fmt.Println("eventSubErr:", err)
	}
	fmt.Printf("%+v\n", resp)
}

func (c *TwitchClient) DeleteSubscriptions(idList []string) {
	for _, id := range idList {
		deleteResp, deleteErr := c.client.RemoveEventSubSubscription(id)

		//TODO integrate logging and handle error
		if deleteErr != nil {
		}
		fmt.Printf("%+v\n", deleteResp)
	}
}

func (c *TwitchClient) GetUsersID(usernameList []string) (idList []string) {
	userResp, userErr := c.client.GetUsers(&helix.UsersParams{
		//IDs:    []string{"twitch user id"},
		//Logins: []string{"twitch user name"},
		Logins: usernameList,
	})

	//TODO integrate logging and handle error
	if userErr != nil {
		fmt.Println("get user error:", userErr)
	}
	fmt.Printf("%+v\n", userResp)
	return idList
}
