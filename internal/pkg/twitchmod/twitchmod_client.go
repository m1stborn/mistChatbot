package twitchmod

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

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

var (
	twitchClientID    = os.Getenv("TWITCH_CLIENT_ID")
	twitchAccessToken = os.Getenv("TWITCH_ACCESSTOKEN")
	secretWord        = "s3cre7w0rd"

	callbackUrl = os.Getenv("CALLBACK_URL")
)

var TC = TwitchClient{}

func init() {
	cl, clErr := helix.NewClient(&helix.Options{
		ClientID:       twitchClientID,
		AppAccessToken: twitchAccessToken,
	})
	//TODO integrate logging and handle error
	if clErr != nil {
		fmt.Println()
	}
	TC.Client = cl
	TC.CallbackUrl = callbackUrl
	TC.SecretWord = secretWord
	logger.Info("init twitch success")
}

func GetSubscriptions() (idList []string) {
	resp, err := TC.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
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

		for _, data := range resp.Data.EventSubSubscriptions {
			idList = append(idList, data.ID)
		}
		logger.WithFields(log.Fields{
			"func": "GetSubscriptions",
		}).Infof(fmt.Sprintf("Subscriptions ID List: %+v", idList))
	}

	return idList //return current eventSubscriptions id
}

func CreateChannelFollowSubscription(broadcasterName string, route string) error {
	idList := GetUsersID([]string{broadcasterName})
	if len(idList) == 0 {
		return errors.New("streamer doesn't exist")
	}
	id := idList[0]

	resp, err := TC.Client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeChannelFollow,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: TC.CallbackUrl + route,
			Secret:   TC.SecretWord,
		},
	})

	if err != nil {
		logger.WithField("func", "CreateChannelFollowSubscription").Error(err.Error())
		return err
	}
	if resp != nil {
		//TODO: logging whole response?
		//logger.WithFields(log.Fields{
		//	"func": "CreateChannelFollowSubscription",
		//}).Infof(fmt.Sprintf("broadcaster ID: %+v, Streamer: %+v", id, broadcasterName))
	}
	return nil
}

func CreateStreamOnlineSubscription(broadcasterName string, route string) error {
	idList := GetUsersID([]string{broadcasterName})
	if len(idList) == 0 {
		return errors.New("streamer doesn't exist")
	}
	id := idList[0]

	resp, err := TC.Client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeStreamOnline,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: id,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: TC.CallbackUrl + route,
			Secret:   TC.SecretWord,
		},
	})

	if err != nil {
		logger.WithField("func", "CreateStreamOnlineSubscription").Error(err.Error())
		return err
	}
	if resp != nil {
		//TODO: logging whole response?
		//logger.WithFields(log.Fields{
		//	"func": "CreateStreamOnlineSubscription",
		//}).Infof(fmt.Sprintf("broadcaster ID: %+v, Streamer: %+v", id, broadcasterName))
	}
	return nil
}

func DeleteSubscriptions(idList []string) {
	for _, id := range idList {
		deleteResp, deleteErr := TC.Client.RemoveEventSubSubscription(id)

		if deleteErr != nil {
			logger.WithField("func", "DeleteSubscriptions").Error(deleteErr.Error())
		}
		if deleteResp != nil {
			//logger.WithFields(log.Fields{
			//	"func": "DeleteSubscriptions",
			//}).Infof(fmt.Sprintf("deleteID: %+v", id))
		}
	}
}

func GetUsersID(usernameList []string) (idList []string) {
	userResp, userErr := TC.Client.GetUsers(&helix.UsersParams{
		//example usage
		//IDs:    []string{"twitch user id"},
		//Logins: []string{"twitch user name"},
		Logins: usernameList,
	})

	if userErr != nil {
		logger.WithField("func", "GetUsersID").Error(userErr.Error())
	}
	if userResp != nil {
		for _, user := range userResp.Data.Users {
			idList = append(idList, user.ID)
		}
		//logger.WithFields(log.Fields{
		//	"func":   "GetUsersID",
		//	"idList": usernameList,
		//}).Infof(fmt.Sprintf("User ID List: %+v", idList))
	}

	return idList
}

func CheckStreamerExist(broadcasterName string) bool {
	idList := GetUsersID([]string{broadcasterName})
	if len(idList) == 0 {
		//return errors.New("streamer doesn't exist"
		return false
	}
	return true
}
