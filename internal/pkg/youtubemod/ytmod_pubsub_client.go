package youtubemod

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Subscription struct {
	topic    string
	id       int
	handler  func(string, []byte) // Content-Type, ResponseBody
	lease    time.Duration
	verified bool
}

func (s Subscription) String() string {
	return fmt.Sprintf("%s (#%d %s)", s.topic, s.id, s.lease)
}

var NIL_SUBSCRIPTION = &Subscription{}

// A HttpRequester is used to make HTTP requests.  http.Client{} satisfies this
// interface.
type HttpRequester interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type PubSubClient struct {
	hubURL string
	self   string

	port          int                      // Which port the server will be started on.
	from          string                   // String passed in the "From" header.
	running       bool                     // Whether the server is running.
	subscriptions map[string]*Subscription // Map of subscriptions.
	httpRequester HttpRequester            // e.g. http.Client{}.
}

func NewPubSubClient(self string, port int, from string) *PubSubClient {
	return &PubSubClient{
		"https://pubsubhubbub.appspot.com/subscribe",
		self,
		port,
		fmt.Sprintf("%s (gohubbub)", from),
		false,
		make(map[string]*Subscription),
		&http.Client{},
	}
}
func (client *PubSubClient) StartClient() {
	client.running = true
}

func (client *PubSubClient) Subscribe(topic string, handler func(string, []byte)) {
	subscription := &Subscription{topic, len(client.subscriptions), handler, 0, false}
	client.subscriptions[topic] = subscription
	if client.running {
		client.makeSubscriptionRequest(subscription)
	}
}

func (client *PubSubClient) Unsubscribe(topic string) {
	if subscription, exists := client.subscriptions[topic]; exists {
		delete(client.subscriptions, topic)
		if client.running {
			client.makeUnsubscribeRequest(subscription)
		}
	} else {
		log.Printf("Cannot unsubscribe, %s doesn't exist.", topic)
	}
}

func (client *PubSubClient) RestoreSubscribe(topic string, id int, handler func(string, []byte)) {
	subscription := &Subscription{topic, id, handler, 0, false}
	client.subscriptions[topic] = subscription
	if client.running {
		client.makeSubscriptionRequest(subscription)
	}
}

func (client PubSubClient) String() string {
	urls := make([]string, len(client.subscriptions))
	i := 0
	for k, _ := range client.subscriptions {
		urls[i] = k
		i++
	}
	return fmt.Sprintf("%d subscription(s): %v", len(client.subscriptions), urls)
}

func (client *PubSubClient) makeSubscriptionRequest(subscription *Subscription) {
	log.Println("Subscribing to", subscription.topic)

	body := url.Values{}
	body.Set("hub.callback", client.formatCallbackURL(subscription.id))
	body.Add("hub.topic", subscription.topic)
	body.Add("hub.mode", "subscribe")
	// body.Add("hub.lease_seconds", "60")

	req, _ := http.NewRequest("POST", client.hubURL, bytes.NewBufferString(body.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("From", client.from)

	resp, err := client.httpRequester.Do(req)

	if err != nil {
		log.Printf("Subscription failed, %s, %s", *subscription, err)

	} else if resp.StatusCode != 202 {
		log.Printf("Subscription failed, %s, status = %s", *subscription, resp.Status)
	}
}

func (client *PubSubClient) makeUnsubscribeRequest(subscription *Subscription) {
	log.Println("Unsubscribing from", subscription.topic)

	body := url.Values{}
	body.Set("hub.callback", client.formatCallbackURL(subscription.id))
	body.Add("hub.topic", subscription.topic)
	body.Add("hub.mode", "unsubscribe")

	req, _ := http.NewRequest("POST", client.hubURL, bytes.NewBufferString(body.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("From", client.from)

	resp, err := client.httpRequester.Do(req)

	if err != nil {
		log.Printf("Unsubscribe failed, %s, %s", *subscription, err)

	} else if resp.StatusCode != 202 {
		log.Printf("Unsubscribe failed, %s status = %d", *subscription, resp.Status)
	}
}

func (client *PubSubClient) formatCallbackURL(callback int) string {
	cbUrl := fmt.Sprintf("https://%s/youtube/pubsub/%d", client.self, callback)
	fmt.Println(cbUrl)
	return cbUrl
}

func (client *PubSubClient) HandlePubSubCallback(resp http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()
	requestBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Printf("Error reading callback request, %s", err)
		return
	}

	params := req.URL.Query()
	topic := params.Get("hub.topic")

	log.Println("HandlePubSubCallback from pub/sub, mode", params.Get("hub.mode"))

	switch params.Get("hub.mode") {
	case "subscribe":
		if subscription, exists := client.subscriptions[topic]; exists {
			subscription.verified = true
			lease, subErr := strconv.Atoi(params.Get("hub.lease_seconds"))
			if subErr == nil {
				subscription.lease = time.Second * time.Duration(lease)
			}
			log.Printf("Subscription verified for %s, lease is %s", topic, subscription.lease)
			resp.Write([]byte(params.Get("hub.challenge")))

		} else {
			log.Printf("Unexpected subscription for %s", topic)
			http.Error(resp, "Unexpected subscription", http.StatusBadRequest)
		}

	case "unsubscribe":
		// We optimistically removed the subscription, so only confirm the
		// unsubscribe if no subscription exists for the topic.
		if _, exists := client.subscriptions[topic]; !exists {
			log.Printf("Unsubscribe confirmed for %s", topic)
			resp.Write([]byte(params.Get("hub.challenge")))

		} else {
			log.Printf("Unexpected unsubscribe for %s", topic)
			http.Error(resp, "Unexpected unsubscribe", http.StatusBadRequest)
		}

	case "denied":
		log.Printf("Subscription denied for %s, reason was %s", topic, params.Get("hub.reason"))
		resp.Write([]byte{})
		// TODO: Don't do anything for now, should probably mark the subscription.

	default:
		subscription, exists := client.subscriptionForPath(req.URL.Path)
		if !exists {
			log.Printf("Callback for unknown subscription: %s", req.URL.String())
			http.Error(resp, "Unknown subscription", http.StatusBadRequest)

		} else {
			log.Printf("Update for %s", subscription)
			resp.Write([]byte{})

			// Asynchronously notify the subscription handler, shouldn't affect response.
			go subscription.handler(req.Header.Get("Content-Type"), requestBody)
		}
	}

}

func (client *PubSubClient) subscriptionForPath(path string) (*Subscription, bool) {
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return NIL_SUBSCRIPTION, false
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		return NIL_SUBSCRIPTION, false
	}
	for _, subscription := range client.subscriptions {
		if subscription.id == id {
			return subscription, true
		}
	}
	return NIL_SUBSCRIPTION, false
}
