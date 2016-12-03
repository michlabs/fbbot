package fbbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Bot struct {
	// User defined fields
	Page            Page // TODO: How to find out it?
	port            int
	verifyToken     string
	pageAccessToken string
	greeting_text   string

	// Handler
	messageHandlers        []MessageHandler
	postbackHandlers       []PostbackHandler
	deliveryHandlers       []DeliveryHandler
	optinHandlers          []OptinHandler
	readHandlers           []ReadHandler
	echoHandlers           []EchoHandler
	checkoutUpdateHandlers []CheckoutUpdateHandler
	paymentHandlers        []PaymentHandler

	// Framework
	Logger *logrus.Logger
	mux    *http.ServeMux
}

type MessageHandler func(*Bot, *Message)
type PostbackHandler func(*Bot, *Postback)
type DeliveryHandler func(*Bot, *Delivery)
type OptinHandler func(*Bot, *Optin)
type ReadHandler func(*Bot, *Read)
type EchoHandler func(*Bot, *Message)
type CheckoutUpdateHandler func(*Bot, *CheckoutUpdate)
type PaymentHandler func(*Bot, *Payment)

func New(port int, verifyToken string, pageAccessToken string) *Bot {
	var b Bot = Bot{
		port:            port,
		verifyToken:     verifyToken,
		pageAccessToken: pageAccessToken,
		mux:             http.NewServeMux(),
		Logger:          logrus.New(),
	}
	b.mux.HandleFunc(WebhookURL, b.handle)
	return &b
}

func (b *Bot) verify(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("hub.mode") == "subscribe" && r.FormValue("hub.verify_token") == b.verifyToken {
		fmt.Fprintf(w, r.FormValue("hub.challenge"))
		b.Logger.Info("Verified")
		return
	}
	b.Logger.WithFields(logrus.Fields{"verifyToken": b.verifyToken, "hub.verify_token": r.FormValue("hub.verify_token")}).Error("Failed to validate. Make sure the validation tokens match.")
	http.Error(w, "Failed validation. Make sure the validation tokens match.", http.StatusForbidden)
	return
}

func (b *Bot) Run() {
	if err := b.Subscribe(); err != nil {
		b.Logger.Fatal("Failed to subscribe to the page")
	}
	if len(b.messageHandlers) == 0 {
		b.Logger.Warn("Message Handler is missing")
	}
	if len(b.postbackHandlers) == 0 {
		b.Logger.Warn("Postback Handler is missing")
	}
	if len(b.deliveryHandlers) == 0 {
		b.Logger.Warn("Delivery Handler is missing")
	}
	if len(b.optinHandlers) == 0 {
		b.Logger.Warn("Optin Handler is missing")
	}
	if len(b.readHandlers) == 0 {
		b.Logger.Warn("Read Handler is missing")
	}
	if len(b.echoHandlers) == 0 {
		b.Logger.Warn("Echo Handler is missing")
	}
	if len(b.checkoutUpdateHandlers) == 0 {
		b.Logger.Warn("Checkout Update Handler is missing")
	}
	if len(b.paymentHandlers) == 0 {
		b.Logger.Warn("Payment Handler is missing")
	}

	b.Logger.Infof("Bot is running at :%d%s", b.port, WebhookURL)
	b.Logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", b.port), b.mux))
}

func (b *Bot) HandleMessage(f MessageHandler) {
	b.messageHandlers = append(b.messageHandlers, f)
}

func (b *Bot) HandlePostback(f PostbackHandler) {
	b.postbackHandlers = append(b.postbackHandlers, f)
}

func (b *Bot) HandleDelivery(f DeliveryHandler) {
	b.deliveryHandlers = append(b.deliveryHandlers, f)
}

func (b *Bot) HandleOptin(f OptinHandler) {
	b.optinHandlers = append(b.optinHandlers, f)
}

func (b *Bot) HandleRead(f ReadHandler) {
	b.readHandlers = append(b.readHandlers, f)
}

func (b *Bot) HandleEcho(f EchoHandler) {
	b.echoHandlers = append(b.echoHandlers, f)
}

func (b *Bot) HandleCheckoutUpdate(f CheckoutUpdateHandler) {
	b.checkoutUpdateHandlers = append(b.checkoutUpdateHandlers, f)
}

func (b *Bot) HandlePayment(f PaymentHandler) {
	b.paymentHandlers = append(b.paymentHandlers, f)
}

func (b *Bot) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		b.verify(w, r)
		return
	}
	if r.Method == "POST" {
		// Handle callback
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			b.Logger.Error("Failed to read resquest body")
			http.Error(w, "Failed to read resquest body", http.StatusInternalServerError)
			return
		}

		b.Logger.WithFields(logrus.Fields{"request": string(body)}).Info("New request")

		var msg rawCallbackMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			b.Logger.WithFields(logrus.Fields{"body": string(body), "error": err.Error()}).Error("Failed to unmarshal request body")
			http.Error(w, "Failed to unmarshal request body", http.StatusInternalServerError)
			return
		}

		// Try to return a 200 OK HTTP as fast as possible
		go b.process(msg.Unbox())

		return
	}
	fmt.Fprintf(w, "Just support GET, POST methods", http.StatusMethodNotAllowed)
}

func (b *Bot) process(messages []interface{}) {
	for _, m := range messages {
		b.Logger.Infof("Message %+v", m)
		switch m := m.(type) {
		case *Message:
			if m.IsEcho {
				for _, f := range b.echoHandlers {
					go f(b, m)
				}
				break
			}
			for _, f := range b.messageHandlers {
				go f(b, m)
			}
		case *Postback:
			for _, f := range b.postbackHandlers {
				go f(b, m)
			}
		case *Delivery:
			for _, f := range b.deliveryHandlers {
				go f(b, m)
			}
		case *Optin:
			for _, f := range b.optinHandlers {
				go f(b, m)
			}
		case *Read:
			for _, f := range b.readHandlers {
				go f(b, m)
			}
		case *CheckoutUpdate:
			for _, f := range b.checkoutUpdateHandlers {
				go f(b, m)
			}
		case *Payment:
			for _, f := range b.paymentHandlers {
				go f(b, m)
			}
		default:
			b.Logger.Error("Unknown message type")
		}
	}
}

// TODO: Refactor
func (b *Bot) Process(message interface{}) {
	var messages []interface{}
	messages = append(messages, message)
	b.process(messages)
}

// TODO: Support other message types
func (b *Bot) Send(r User, m interface{}) error {
	switch m := m.(type) {
	case *TextMessage:
		return b.sendTextMessage(r, m)
	case *ImageMessage:
		return b.sendImageMessage(r, m)
	case *ButtonMessage:
		return b.sendButtonMessage(r, m)
	case *GenericMessage:
		return b.sendGenericMessage(r, m)
	default:
		return errors.New("unknown message type")
	}
}

func (b *Bot) SendText(r User, text string) error {
	m := NewTextMessage(text)
	return b.Send(r, m)
}

func (b *Bot) sendTextMessage(r User, m *TextMessage) error {
	data := make(map[string]interface{})
	data["notification_type"] = m.Noti
	data["recipient"] = map[string]string{"id": r.ID}
	data["message"] = map[string]string{"text": m.Text}

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
		return err
	}

	return nil
}

func (b *Bot) sendImageMessage(r User, m *ImageMessage) error {
	payload := make(map[string]string)
	payload["url"] = m.URL

	attachment := make(map[string]interface{})
	attachment["type"] = "image"
	attachment["payload"] = payload

	message := make(map[string]interface{})
	message["attachment"] = attachment

	data := make(map[string]interface{})
	data["recipient"] = r
	data["message"] = message
	data["notification_type"] = m.Noti

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
	}

	return nil
}

func (b *Bot) sendButtonMessage(r User, m *ButtonMessage) error {
	data := make(map[string]interface{})

	payload := make(map[string]interface{})
	payload["template_type"] = "button"
	payload["text"] = m.Text
	payload["buttons"] = m.Buttons

	attachment := make(map[string]interface{})
	attachment["type"] = "template"
	attachment["payload"] = payload

	data["notification_type"] = m.Noti
	data["recipient"] = map[string]string{"id": r.ID}
	data["message"] = map[string]interface{}{"attachment": attachment}

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
		return err
	}

	return nil
}

func (b *Bot) sendGenericMessage(r User, m *GenericMessage) error {
	payload := make(map[string]interface{})
	payload["template_type"] = "generic"
	payload["elements"] = m.Bubbles

	attachment := make(map[string]interface{})
	attachment["type"] = "template"
	attachment["payload"] = payload

	data := make(map[string]interface{})
	data["notification_type"] = m.Noti 
	data["recipient"] = map[string]string{"id": r.ID}
	data["message"] = map[string]interface{}{"attachment": attachment}

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
		return err
	}

	return nil
}

func (b *Bot) TypingOn(r User) error {
	data := make(map[string]interface{})
	data["recipient"] = r
	data["sender_action"] = "typing_on"

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
		return err
	}

	return nil
}

func (b *Bot) TypingOff(r User) error {
	data := make(map[string]interface{})
	data["recipient"] = r
	data["sender_action"] = "typing_off"

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
		return err
	}

	return nil
}

func (b *Bot) MarkSeen(r User) error {
	data := make(map[string]interface{})
	data["recipient"] = r
	data["sender_action"] = "mark_seen"

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data, "error": err}).Error("Failed to send message")
		return err
	}

	return nil
}

// Subscribe subscribes this bot to get updates for the page.
func (b *Bot) Subscribe() error {
	data := make(map[string]interface{})
	if resp, err := b.httppost(APIEndpoint+"/me/subscribed_apps", data); err != nil {
		b.Logger.WithFields(logrus.Fields{"error": err, "resp": resp}).Error("Failed to subscribe")
		return err
	}
	b.Logger.Info("Success to subscribe your page")
	return nil
}

// TODO
func (b *Bot) AddGreetingText(text string) error {
	return nil
	// greeting.text must be UTF-8 and has a 160 character limit
	// curl -X POST -H "Content-Type: application/json" -d '{
	//   "setting_type":"greeting",
	//   "greeting":{
	//     "text":"Timeless apparel for the masses."
	//   }
	// }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=PAGE_ACCESS_TOKEN"

	// curl -X POST -H "Content-Type: application/json" -d '{
	//   "setting_type":"greeting",
	//   "greeting":{
	//     "text":"Hi {{user_first_name}}, welcome to this bot."
	//   }
	// }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=PAGE_ACCESS_TOKEN"
}

// TODO
func (b *Bot) RemoveGreetingText() error {
	return nil
	// curl -X DELETE -H "Content-Type: application/json" -d '{
	//   "setting_type":"greeting"
	// }' "https://graph.facebook.com/v2.6/me/thread_settings?access_token=PAGE_ACCESS_TOKEN"
}

func (b *Bot) httppost(url string, data map[string]interface{}) ([]byte, error) {
	url = fmt.Sprintf("%s?access_token=%s", url, b.pageAccessToken)

	d, err := json.Marshal(data)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"data": data}).Error("Failed to marshal")
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"URL": url, "data": data}).Error("Failed to request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Logger.WithFields(logrus.Fields{"body": resp.Body}).Error("Failed to read response body")
		return nil, err
	}

	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Response code: %d. Body: %s", resp.StatusCode, string(body)))
		b.Logger.WithFields(logrus.Fields{"error": err.Error()}).Error("Request is not success")
		return nil, err
	}

	return body, nil
}
