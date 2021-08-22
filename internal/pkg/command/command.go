package command

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/m1stborn/mistChatbot/internal/pkg/model"
	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"
)

//func HandelLineFollow(id, accountType string) {
//	//fmt.Println(model.TestUser)
//}

var streamOnlineRoute = "/twitch/streamOnline"

func HandleCommand(text string, user *model.User, isUser bool) string {
	command := strings.ToLower(strings.Fields(strings.TrimSpace(text))[0])
	if isUser {
		//TODO Handel log
	}
	switch command {
	case "/sub":
		//step 0: regex match command
		re := regexp.MustCompile("^(/sub)\\s([a-zA-Z0-9_]{4,25}$)")
		if matched := re.MatchString(text); !matched {
			//TODO handle reply message
			return "wrong format of command"
		}
		args := re.FindStringSubmatch(text)
		streamerName := args[2]

		//step 1: check if the broadcaster exist
		if !twitchmod.CheckStreamerExist(streamerName) {
			return "streamer not exist"
		}
		//step 2: check if already sub to twitch EventSub
		if !model.DB.CheckStreamerExist(streamerName) {
			err := twitchmod.CreateStreamOnlineSubscription(streamerName, streamOnlineRoute)
			if err != nil {
				//TODO handle error
			}
		}
		//step 3: write into DB
		model.DB.CreateSubscription(&model.Subscription{
			Line:            user.Line,
			LineAccessToken: user.LineAccessToken,
			TwitchLoginName: streamerName,
		})
		logger.WithFields(log.Fields{
			"pkg":  "command",
			"case": "/sub",
			"func": "HandleCommand",
		}).Info("sub success")
		return fmt.Sprintf("sub %v successful!", streamerName)
	}
	//case "/del":
	return "no this command, please check /help"
}
