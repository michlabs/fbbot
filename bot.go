package fbbot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/michlabs/fbbot/memory"
	"io/ioutil"
	"net/http"
)

type Bot struct {
	// User defined fields
	Page            Page // TODO: How to find out it?
	port            int
	verifyToken     string
	appSecret       string
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

	LTMemory memory.Memory // LTMemory will be persit across conversation
	STMemory memory.Memory // STMemory will be cleared for the user at the end of conversation

	// Framework
	Logger *logrus.Logger
	mux    *http.ServeMux
}

func New(port int, verifyToken string, appSecret string, pageAccessToken string) *Bot {
	var b Bot = Bot{
		port:            port,
		verifyToken:     verifyToken,
		appSecret:       appSecret,
		pageAccessToken: pageAccessToken,
		mux:             http.NewServeMux(),
		Logger:          logrus.New(),
	}
	b.mux.HandleFunc(WebhookURL, b.handle)
	b.LTMemory = memory.New("ephemeral")
	b.STMemory = memory.New("ephemeral")
	bot = &b // For using outside of bot methods (User struct)
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

func (b *Bot) AddMessageHandler(h MessageHandler) {
	b.messageHandlers = append(b.messageHandlers, h)
}

func (b *Bot) AddPostbackHandler(h PostbackHandler) {
	b.postbackHandlers = append(b.postbackHandlers, h)
}

func (b *Bot) AddDeliveryHandler(h DeliveryHandler) {
	b.deliveryHandlers = append(b.deliveryHandlers, h)
}

func (b *Bot) AddOptinHandler(h OptinHandler) {
	b.optinHandlers = append(b.optinHandlers, h)
}

func (b *Bot) AddReadHandler(h ReadHandler) {
	b.readHandlers = append(b.readHandlers, h)
}

func (b *Bot) AddEchoHandler(h EchoHandler) {
	b.echoHandlers = append(b.echoHandlers, h)
}

func (b *Bot) AddCheckoutUpdateHandler(h CheckoutUpdateHandler) {
	b.checkoutUpdateHandlers = append(b.checkoutUpdateHandlers, h)
}

func (b *Bot) AddPaymentHandler(h PaymentHandler) {
	b.paymentHandlers = append(b.paymentHandlers, h)
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
		b.Logger.WithFields(logrus.Fields{"request": string(body)}).Debug("New request:")

		// Verify message signature
		if !b.verifySignature(body, r.Header.Get("X-Hub-Signature")[5:]) {
			b.Logger.Error("invalid request signature")
			return
		}

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

	http.Error(w, "Just support GET, POST methods", http.StatusMethodNotAllowed)
}

func (b *Bot) process(messages []interface{}) {
	for _, m := range messages {
		b.Logger.Debugf("Message %+v", m)
		switch m := m.(type) {
		case *Message:
			if m.IsEcho {
				for _, h := range b.echoHandlers {
					go h.HandleEcho(b, m)
				}
				break
			}
			for _, h := range b.messageHandlers {
				go h.HandleMessage(b, m)
			}
		case *Postback:
			for _, h := range b.postbackHandlers {
				go h.HandlePostback(b, m)
			}
		case *Delivery:
			for _, h := range b.deliveryHandlers {
				go h.HandleDelivery(b, m)
			}
		case *Optin:
			for _, h := range b.optinHandlers {
				go h.HandleOptin(b, m)
			}
		case *Read:
			for _, h := range b.readHandlers {
				go h.HandleRead(b, m)
			}
		case *CheckoutUpdate:
			for _, h := range b.checkoutUpdateHandlers {
				go h.HandleCheckoutUpdate(b, m)
			}
		case *Payment:
			for _, h := range b.paymentHandlers {
				go h.HandlePayment(b, m)
			}
		default:
			b.Logger.Error("Unknown message type")
		}
	}
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
	case *QuickRepliesMessage:
		return b.sendQuickRepliesMessage(r, m)
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
	data["messaging_type"] = "RESPONSE"
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

// SendImage sends an image specified by an URL to the receipient
func (b *Bot) SendImage(r User, url string) error {
	m := NewImageMessage()
	m.URL = url
	return b.Send(r, m)
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
	data["messaging_type"] = "RESPONSE"
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

	data["messaging_type"] = "RESPONSE"
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
	data["messaging_type"] = "RESPONSE"
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

func (b *Bot) sendQuickRepliesMessage(r User, m *QuickRepliesMessage) error {
	data := make(map[string]interface{})
	data["messaging_type"] = "RESPONSE"
	data["recipient"] = map[string]string{"id": r.ID}
	data["message"] = m

	_, err := b.httppost(SendAPIEndpoint, data)
	if err != nil {
		b.Logger.Errorf("Failed to send message. Error: %s\nData:%#v", err.Error(), data)
		return err
	}

	return nil
}

func (b *Bot) TypingOn(r User) error {
	data := make(map[string]interface{})
	data["messaging_type"] = "RESPONSE"
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
	data["messaging_type"] = "RESPONSE"
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
	data["subscribed_fields"] = []string{
		"message_mention", 
		"messages", 
		"message_reactions", 
		"messaging_account_linking", 
		"messaging_checkout_updates",
		"message_echoes", 
		"message_deliveries", 
		"messaging_optins", 
		"messaging_optouts", 
		"messaging_payments", 
		"messaging_postbacks", 
		"messaging_pre_checkouts", 
		"message_reads", 
		"messaging_referrals", 
		"messaging_handovers", 
		"messaging_policy_enforcement", 
		"messaging_page_feedback", 
		"messaging_appointments", 
		"messaging_direct_sends", 
		"messaging_fblogin_account_linking"
		"messaging_feedback",
	}
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

func (b *Bot) fetchUserData(u *User) {
	uri := fmt.Sprintf("%s/%s?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%s", APIEndpoint, u.ID, b.pageAccessToken)
	b.Logger.Debug("fetchUserData", uri)

	resp, err := http.Get(uri)
	if err != nil {
		b.Logger.Error("failed to fetch user data: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Logger.Error("failed to read response data")
		return
	}

	if resp.StatusCode != 200 {
		b.Logger.Error("failed to fetch user data: ", string(body))
		return
	}

	var tmp struct {
		FirstName        string  `json:"first_name, omitempty"`
		LastName         string  `json:"last_name, omitempty"`
		ProfilePic       string  `json:"profile_pic, omitempty"`
		Locale           string  `json:"locale, omitempty"`
		Timezone         float32 `json:"timezone, omitempty"`
		Gender           string  `json:"gender, omitempty"`
		IsPaymentEnabled bool    `json:"is_payment_enabled, omitempty"` // Is the user eligible to receive messenger platform payment messages
	}

	if err := json.Unmarshal(body, &tmp); err != nil {
		b.Logger.Error("failed to unmarshal response: ", err)
		return
	}

	u.firstName = tmp.FirstName
	u.lastName = tmp.LastName
	u.profilePic = tmp.ProfilePic
	u.locale = tmp.Locale
	u.timezone = tmp.Timezone
	u.gender = tmp.Gender
	u.isPaymentEnabled = tmp.IsPaymentEnabled
	u.isFetched = true

	return
}

// EnableGetStarted enables the Get Started button at the first conversation
// payload will be sent back to the bot when user clicks on the button.
func (b *Bot) EnableGetStarted(payload string) error {
	getstarted := make(map[string]string)
	getstarted["payload"] = payload

	data := make(map[string]interface{})
	data["get_started"] = getstarted
	_, err := b.httppost(ProfileEndpoint, data)
	return err
}

func (b *Bot) AddPersistentMenus(menus ...*Menu) error {
	data := make(map[string]interface{})
	data["persistent_menu"] = menus
	_, err := b.httppost(ProfileEndpoint, data)
	return err
}

func (b *Bot) verifySignature(content []byte, signature string) bool {
	if signature == "" {
		return false
	}
	mac := hmac.New(sha1.New, []byte(b.appSecret))
	mac.Write(content)
	if fmt.Sprintf("%x", mac.Sum(nil)) != signature {
		return false
	}
	return true
}
