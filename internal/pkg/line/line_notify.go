package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ht "html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"

	"github.com/m1stborn/mistChatbot/internal/pkg/model"

	log "github.com/sirupsen/logrus"
)

const notifyBotHost string = "https://notify-bot.line.me"

var (
	params       map[string]string
	clientID     = os.Getenv("LINE_NOTIFY_CLIENT_ID")
	clientSecret = os.Getenv("LINE_NOTIFY_CLIENT_SECRET")
	redirectURI  = os.Getenv("CALLBACK_URL") + "/line/notify/callback"
)

func buildQueryString(params map[string]string) (query string) {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		query += fmt.Sprintf("%s=%s&", key, params[key])
	}
	return query
}

func getAuthorizeURL(lineID string) string {
	var uri = "/oauth/authorize"
	params = map[string]string{
		"response_type": "code",
		"client_id":     clientID,
		"redirect_uri":  redirectURI,
		"scope":         "notify",
		"state":         lineID,
		"response_mode": "form_post",
	}
	query := buildQueryString(params)
	return fmt.Sprintf("%s%s?%s", notifyBotHost, uri, query)
}

//TODO: serve a frontend website
func fetchAccessToken(code string) (string, error) {
	type responseBody struct {
		AccessToken string `json:"access_token"`
	}
	uri := "/oauth/token"
	params = map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectURI,
		"client_id":     clientID,
		"client_secret": clientSecret,
	}
	body := buildQueryString(params)
	r, httpErr := http.Post(notifyBotHost+uri, "application/x-www-form-urlencoded", bytes.NewBufferString(body))
	if httpErr != nil {
		logger.WithField("func", "fetchAccessToken").
			WithError(httpErr).Error("Post Error")
	}
	if r.StatusCode != http.StatusOK {
		lineErr := errors.New("Get Line Access Token Error, StatusCode:" + strconv.Itoa(r.StatusCode))
		logger.WithField("func", "fetchAccessToken").
			WithError(lineErr).Error()
		return "", err
	}
	var rspBody responseBody
	decodeErr := json.NewDecoder(r.Body).Decode(&rspBody)
	if decodeErr != nil {
		logger.WithField("func", "fetchAccessToken").
			WithError(decodeErr).Error("Decode Line Access Token Error")
		return "", decodeErr
	}
	return rspBody.AccessToken, nil
}

//TODO legacy , now directly post url

func HandelNotifyAuth(w http.ResponseWriter, r *http.Request) {
	t, tErr := ht.New("webpage").Parse(authTmpl)
	if tErr != nil {
		logger.WithFields(log.Fields{
			"func": "HandelNotifyAuth",
		}).Error("Create template error!")
	}
	noItems := struct {
		ClientID    string
		CallbackURL string
	}{
		ClientID:    clientID,
		CallbackURL: callbackUrl + "/line/notify/callback",
	}
	tErr = t.Execute(w, noItems)
	if tErr != nil {
		logger.WithFields(log.Fields{
			"func": "HandelNotifyAuth",
		}).Error("Execute template error!")
	}
}

func HandleNotifyCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("error") != "" {
		logger.WithFields(log.Fields{
			"error":       r.FormValue("error"),
			"state":       r.FormValue("state"),
			"description": r.FormValue("error_description"),
		}).Info("Get LINE Notify Callback Failed")
	}

	code, lineID := r.FormValue("code"), r.FormValue("state")
	accessToken, tokenErr := fetchAccessToken(code)
	if tokenErr != nil {
		logger.WithError(tokenErr).Info("Fetch Access Token Failed")
	}

	dbErr := model.DB.UserConnectNotify(lineID, accessToken)
	if dbErr != nil {
		SendLineNotify(accessToken, "連結 LINE Notify 失敗")
	} else {
		SendLineNotify(accessToken, "連結 LINE Notify 成功, 請回到 mistChatbot 輸入 /help 查看相關功能!")
		//if _, err = bot.PushMessage(lineID, linebot.NewTextMessage("連結 LINE Notify 成功, 輸入 /help 查看相關功能!")).Do(); err != nil {
		//	logger.WithFields(log.Fields{
		//		"func":   "HandleNotifyCallback",
		//		"lineID": lineID,
		//	}).Error(err)
		//}
	}

	logger.WithFields(log.Fields{
		"func":        "HandleNotifyCallback",
		"lineID":      lineID,
		"accessToken": accessToken,
	}).Info("Successfully register user!")
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
