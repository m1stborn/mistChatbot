package youtubemod

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	//"github.com/go-playground/validator/v10"
)

var (
	YtToken = os.Getenv("YOUTUBE_DATA_API_TOKEN")

	dataApiBaseUrl = "https://youtube.googleapis.com/youtube/v3/videos"
)

// VideoResource TODO: implement pagination
type VideoResource struct {
	Kind          string   `json:"kind"`
	Etag          string   `json:"etag"`
	NextPageToken string   `json:"nextPageToken"`
	PrevPageToken string   `json:"prevPageToken"`
	PageInfo      PageInfo `json:"pageInfo"`
	Items         []Video  `json:"items"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Video struct {
	Id                   string                `json:"id" validate:"required"`
	Snippet              Snippet               `json:"snippet" validate:"required"`
	LiveStreamingDetails *LiveStreamingDetails `json:"liveStreamingDetails"`
}

type LiveStreamingDetails struct {
	ActualStartTime    *time.Time `json:"actualStartTime"`
	ActualEndTime      *time.Time `json:"actualEndTime"`
	ScheduledStartTime time.Time  `json:"scheduledStartTime"`
	ScheduledEndTime   time.Time  `json:"scheduledEndTime"`
	ConcurrentViewers  string     `json:"concurrentViewers"`
	ActiveLiveChatId   string     `json:"activeLiveChatId"`
}

type Snippet struct {
	ChannelID   string `json:"channelId" validate:"required"`
	Title       string `json:"title" validate:"required"`
	ChannelName string `json:"channelTitle" validate:"required"`
}

type VideoItems struct {
	Items []Video `json:"items" validate:"dive"`
}

func GetStreamIdLiveDetailByIds(vidIds []string) (VideoItems, error) {
	//dataApiBaseUrl := "https://youtube.googleapis.com/youtube/v3/videos"

	u, urlErr := url.Parse(dataApiBaseUrl)
	if urlErr != nil {
		fmt.Println(urlErr)
	}

	params := url.Values{}

	ids := strings.Join(vidIds, ",")

	params.Add("part", "liveStreamingDetails,snippet")
	params.Add("fields", "items(id,snippet(channelId,title,channelTitle),liveStreamingDetails)")
	params.Add("id", ids)
	params.Add("key", YtToken)

	u.RawQuery = params.Encode()

	ytResp, urlErr := http.Get(u.String())
	if urlErr != nil {
		// handle error
		fmt.Println("urlErr", urlErr)
	}

	defer ytResp.Body.Close()
	body, urlErr := ioutil.ReadAll(ytResp.Body)
	if urlErr != nil {
		// handle error
		fmt.Println("urlErr", urlErr)
	}

	//var jsonMsg map[string]interface{}
	//var videoObjs VideoResource
	var videoObjs VideoItems

	jsonErr := json.Unmarshal(body, &videoObjs)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}

	//for _, video := range videoObjs.Items {
	//	if video.LiveStreamingDetails != nil {
	//		if video.LiveStreamingDetails.ActualStartTime == nil {
	//			fmt.Printf("%s stream not start yet\n", video.Snippet.Title)
	//		}
	//	}
	//}

	//Validate the response not missing.
	validate := validator.New()
	validErr := validate.Struct(videoObjs)
	if validErr != nil {
		return videoObjs, validErr
	}
	return videoObjs, nil
}
