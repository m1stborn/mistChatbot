package line

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	ht "html/template"
	"net/http"
	"os"
	"sort"
	"strconv"

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
		}).Error("Get LINE Notify Callback Failed")
	}

	code, lineID := r.FormValue("code"), r.FormValue("state")
	accessToken, tokenErr := fetchAccessToken(code)
	if tokenErr != nil {
		logger.WithError(tokenErr).Error("Fetch Access Token Failed")
	}

	logger.WithFields(log.Fields{
		"func":        "CatchCallback",
		"lineID":      lineID,
		"accessToken": accessToken,
	}).Info("Successfully register user!")
}
