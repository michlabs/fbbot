package fbbot

var bot *Bot // For using outside bot's method, for example: in User struct

const (
	WebhookURL      = "/webhook"
	SendAPIEndpoint = "https://graph.facebook.com/v2.6/me/messages"
	APIEndpoint     = "https://graph.facebook.com/v2.6"
	ProfileEndpoint = "https://graph.facebook.com/v2.6/me/messenger_profile"

	// Notification type
	NotiRegular    string = "REGULAR"     // will emit a sound/vibration and a phone notification
	NotiSilentPush string = "SILENT_PUSH" // will just emit a phone notification
	NotiNoPush     string = "NO_PUSH"     // will not emit either
)
