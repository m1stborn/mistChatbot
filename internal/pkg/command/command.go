package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/m1stborn/mistChatbot/internal/pkg/model"
	"github.com/m1stborn/mistChatbot/internal/pkg/twitchmod"
)

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
			return "Wrong format of command"
		}
		args := re.FindStringSubmatch(text)
		streamerName := args[2]

		//step 1: check if the broadcaster exist
		if !twitchmod.CheckStreamerExist(streamerName) {
			return "Streamer not exist"
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
			//TODO TwitchEventSubID
		})
		return fmt.Sprintf("sub %v successful!", streamerName)
	case "/del":
		//step 0: regex match command
		re := regexp.MustCompile("^(/del)\\s([a-zA-Z0-9_]{4,25}$)")
		if matched := re.MatchString(text); !matched {
			//TODO handle reply message
			return "Wrong format of command"
		}
		args := re.FindStringSubmatch(text)
		streamerName := args[2]
		//step 1: delete DB record
		err := model.DB.DeleteSubByUserBroadcaster(user.Line, streamerName)
		if err == model.ErrRecordNotExist {
			return "Wrong streamer name or not sub yet "
		}
		//step 2: check if should remove eventSub from twitch
		if !model.DB.CheckStreamerExist(streamerName) {
			//TODO TwitchEventSubID
		}
		return fmt.Sprintf("delete %v successful!", streamerName)
	case "/list":
		subs := model.DB.QuerySubByUser(user.Line)
		if len(subs) == 0 {
			return "You haven't subscribe any channel!"
		}
		resp := "Currently subs:\n"
		for _, sub := range subs {
			resp += fmt.Sprintf("https://www.twitch.tv/%s\n", sub.TwitchLoginName)
		}
		return resp
	case "/help":
		//TODO add  line emoji
		return fmt.Sprint("指令清單:\n" +
			"1. /sub [twitch ID]: 訂閱頻道\n   Example: /sub never_loses\n" +
			"2. /del [twitch ID]: 刪除頻道\n   Example: /del qq7925168\n" +
			"3. /list: 列出訂閱的頻道")
	}
	return "No this command, please check /help"
}
