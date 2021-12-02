package youtubemod

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/m1stborn/mistChatbot/internal/pkg/model"

	log "github.com/sirupsen/logrus"
)

type YtTracker struct {
	Upcoming map[string]*YtStream
	VideoIds []string
}

type YtStream struct {
	ChannelId   string
	ChannelName string
	Title       string
	VideoId     string
	VideoUrl    string
	//LiveStreamingDetails LiveStreamingDetails // actually need?
}

var Tracker = YtTracker{}

var (
	videoBaseUrl = "https://www.youtube.com/watch?v="
)

func (y *YtTracker) Init() {
	y.Upcoming = make(map[string]*YtStream, 0)
}

func (y *YtTracker) StartTrack() {
	// Restore VideoIDs from DB
	var oldVideos []string
	oldVideos = model.DB.QueryAllVideoIds()
	oldVideoJson, testErr := GetStreamIdLiveDetailByIds(oldVideos)
	if testErr != nil {
		fmt.Println("DataApi miss field:", testErr)
		fmt.Printf("video_resource:%+v\n", oldVideoJson)
	} else {
		y.Update(oldVideoJson)
	}
	//Implementation 1
	//ticker := time.NewTicker(5 * time.Second)
	//quit := make(chan struct{})
	//go func() {
	//	for {
	//		select {
	//		case <- ticker.C:
	//			// do stuff
	//		case <- quit:
	//			ticker.Stop()
	//			return
	//		}
	//	}
	//}()

	//Implementation 2
	tick := time.Tick(60 * time.Second)
	for range tick {
		//Send YouTube Data api request to check if stream start.
		if len(y.VideoIds) > 0 {
			fmt.Println("videoIds:", y.VideoIds)
		}

		videoJson, err := GetStreamIdLiveDetailByIds(y.VideoIds)
		if err != nil {
			fmt.Println("DataApi miss field:", err)
			fmt.Printf("video_resource:%+v\n", videoJson)
			continue
		}
		y.Update(videoJson)
	}

}

func (y *YtTracker) Update(videoJson VideoItems) {
	for _, video := range videoJson.Items {

		//fmt.Printf("video:%+v\n", video)

		//step 1: This video is not a stream, discard and remove from Upcoming list
		if video.LiveStreamingDetails == nil {
			delete(y.Upcoming, video.Id)
			y.VideoIds = remove(y.VideoIds, video.Id)
			model.DB.DeleteYtVideo(video.Id)
			continue
		}

		//step 2: Check if the video is on the upcoming list
		if _, ok := y.Upcoming[video.Id]; !ok {
			y.Upcoming[video.Id] = &YtStream{
				VideoId:     video.Id,
				Title:       video.Snippet.Title,
				ChannelId:   video.Snippet.ChannelID,
				ChannelName: video.Snippet.ChannelName,
				VideoUrl:    videoBaseUrl + video.Id,
			}
			y.VideoIds = append(y.VideoIds, video.Id)
		}

		//step 3: Check if ChannelName is empty
		if y.Upcoming[video.Id].ChannelName == "" {
			//The ChannelName == "" means it just been added to Upcoming list.
			y.Upcoming[video.Id].ChannelName = video.Snippet.ChannelName
		}

		//step 4: The upcoming video start streaming, send out the Notification and remove from Upcoming list
		if video.LiveStreamingDetails.ActualStartTime != nil {
			if video.LiveStreamingDetails.ActualEndTime == nil {
				//TODO: parallel run
				//TODO: For testing comment out this loop (due to no Database!!!)
				accessTokens := model.DB.QuerySubByYtChannelId(video.Snippet.ChannelID)
				//
				for _, token := range accessTokens {
					SendLineNotify(token,
						fmt.Sprintf("%s start streaming!\n %s",
							video.Snippet.ChannelName,
							y.Upcoming[video.Id].VideoUrl,
						))
				}
				fmt.Printf("%s start streaming!\n%s\n",
					video.Snippet.ChannelName,
					y.Upcoming[video.Id].VideoUrl,
				)
			}

			delete(y.Upcoming, video.Id)
			y.VideoIds = remove(y.VideoIds, video.Id)
			model.DB.DeleteYtVideo(video.Id)
		}
	}
}

func (y *YtTracker) AddUpcoming(feed Feed) {
	for _, entry := range feed.Entries {
		//If the VideoId already exist, means there is some update to the Title
		if _, ok := y.Upcoming[entry.VideoID]; ok {
			y.Upcoming[entry.VideoID].Title = feed.Title
			continue
		}
		//The Feed won't have liveStreamingDetail, should be updated soon
		y.Upcoming[entry.VideoID] = &YtStream{
			VideoId:     entry.VideoID,
			Title:       entry.Title,
			ChannelId:   entry.ChannelId,
			ChannelName: "",
			VideoUrl:    videoBaseUrl + entry.VideoID,
		}
		y.VideoIds = append(y.VideoIds, entry.VideoID)
		model.DB.CreateYtVideo(&model.YtVideo{
			VideoId: entry.VideoID,
		})
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
		log.Println("Line Notify Failed", data)
		//logger.WithFields(log.Fields{
		//	"status":   r.Status,
		//	"response": string(data),
		//}).Error("LINE Notify Failed")
	}
}
