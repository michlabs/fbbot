package fbbot

import (
	"github.com/Sirupsen/logrus"
)

// rawCallbackMessage is data you will receive at your webhook
type rawCallbackMessage struct {
	// Object indicates the object type that this payload applies to.
	// It could be: user, page, permissions, payments.
	// In this case, value will alway be "page".
	RawObject string `json:"object"`

	// rawEntries is a slice containing event data of the same object type
	// that are batched together
	RawEntries []struct {
		// rawID is ID of the object that triggers this event.
		// In this case, it will alway be Page ID
		RawID string `json:"id"`

		// rawEventTime is time the event data was sent
		RawEventTime int64 `json:"time"`

		// rawMessaging contains data related to messaging
		RawMessaging []rawMessageData `json:"messaging"`
	} `json:"entry"`
}

// rawMessageData contains data related to a message
type rawMessageData struct {
	// rawSender is user that sends message
	RawSender User `json:"sender"`

	// rawRecipient is user that receives message
	// In this case, RawRecipient is your Page
	RawRecipient Page `json:"recipient"`

	// rawTimestamp is the time messageData was sent
	RawTimestamp int64 `json:"timestamp"`

	RawMessage     *rawMessage     `json:"message"`
	Postback       *Postback       `json:"postback"`
	Delivery       *Delivery       `json:"delivery"`
	Optin          *Optin          `json:"optin"`
	Read           *Read           `json:"read"`
	CheckoutUpdate *CheckoutUpdate `json:"checkout_update"`
	Payment        *Payment        `json:"payment"`
}

// rawMessage is a Facebook message
type rawMessage struct {
	// rawMid is message ID
	RawMid string `json:"mid"`

	// rawSeq is message sequence number
	RawSeq int `json:"seq"`

	// rawText is text of message
	RawText string `json:"text"`

	RawIsEcho bool `json:"is_echo"`

	RawAppID int64 `json:"app_id"`

	// rawAttachments is a slice containing attachment data
	RawAttachments []rawAttachment `json:"attachments"`
}

// rawAttachment is attached image, video, audio or location
type rawAttachment struct {
	// rawType is type of the attachment
	// It could be: image, video, audio or location
	RawType string `json:"type"`

	// Attachment file
	RawPayload rawPayload `json:"payload"`
}

// Attachment file
type rawPayload struct {
	// URL of the attachment file
	RawURL         string         `json:"url"`
	RawCoordinates rawCoordinates `json:"coordinates"`
}

type rawCoordinates struct {
	RawLat  float64 `json:"lat"`
	RawLong float64 `json:"long"`
}

func (cbMsg *rawCallbackMessage) Unbox() []interface{} {
	var messages []interface{}
	for _, entry := range cbMsg.RawEntries {
		for _, rawMessageData := range entry.RawMessaging {
			if rawMessageData.RawMessage != nil {
				messages = append(messages, buildMessage(rawMessageData))
			} else if rawMessageData.Postback != nil {
				rawMessageData.Postback.Sender = rawMessageData.RawSender
				messages = append(messages, rawMessageData.Postback)
			} else if rawMessageData.Delivery != nil {
				messages = append(messages, rawMessageData.Delivery)
			} else if rawMessageData.Optin != nil {
				rawMessageData.Optin.Sender = rawMessageData.RawSender
				messages = append(messages, rawMessageData.Optin)
			} else if rawMessageData.Read != nil {
				rawMessageData.Read.Sender = rawMessageData.RawSender
				messages = append(messages, rawMessageData.Read)
			} else if rawMessageData.CheckoutUpdate != nil {
				rawMessageData.CheckoutUpdate.Sender = rawMessageData.RawSender
				messages = append(messages, rawMessageData.CheckoutUpdate)
			} else if rawMessageData.Payment != nil {
				rawMessageData.Payment.Sender = rawMessageData.RawSender
				messages = append(messages, rawMessageData.Payment)
			} else {
				logrus.WithFields(logrus.Fields{"rawMessageData": rawMessageData}).Error("Unknown message type")
			}
		}
	}
	return messages
}

func buildMessage(m rawMessageData) *Message {
	var msg Message
	msg.ID = m.RawMessage.RawMid
	msg.Page = m.RawRecipient
	msg.Sender = m.RawSender
	msg.Text = m.RawMessage.RawText
	msg.Seq = m.RawMessage.RawSeq
	msg.Timestamp = m.RawTimestamp
	msg.IsEcho = m.RawMessage.RawIsEcho
	msg.AppID = m.RawMessage.RawAppID
	for _, attachment := range m.RawMessage.RawAttachments {
		switch attachment.RawType {
		case "image":
			image := Image{URL: attachment.RawPayload.RawURL}
			msg.Images = append(msg.Images, image)
		case "video":
			video := Video{URL: attachment.RawPayload.RawURL}
			msg.Videos = append(msg.Videos, video)
		case "audio":
			audio := Audio{URL: attachment.RawPayload.RawURL}
			msg.Audios = append(msg.Audios, audio)
		case "file":
			file := File{URL: attachment.RawPayload.RawURL}
			msg.Files = append(msg.Files, file)
		case "location":
			location := Location{
				Coordinates: Coordinates{
					Lat:  attachment.RawPayload.RawCoordinates.RawLat,
					Long: attachment.RawPayload.RawCoordinates.RawLong,
				},
			}
			msg.Location = location
		}

	}
	return &msg
}
