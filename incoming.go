package fbbot

// Message represents a message sent to your page.
type Message struct {
	ID         string
	Page       Page
	Sender     User
	Text       string
	IsEcho     bool
	AppID      int64
	Images     []Image
	Videos     []Video
	Audios     []Audio
	Files      []File
	Location   Location
	Seq        int
	Timestamp  int64
	Quickreply Quickreply
}

type Quickreply struct {
	Payload string
}

type Image struct {
	URL string
}

type Video struct {
	URL string
}

type Audio struct {
	URL string
}

type File struct {
	URL string
}

type Location struct {
	Coordinates Coordinates
}

type Coordinates struct {
	Lat  float64
	Long float64
}

// Postback
type Postback struct {
	Sender  User
	Payload string `json:"payload"`
}

// Delivery
// This callback will occur when a message a page has sent has been delivered.
type Delivery struct {
	MessageIDs []string `json:"mids"`      // Slice containing message IDs of messages that were delivered. Field may not be present.
	Watermark  float64  `json:"watermark"` // All messages that were sent before this timestamp were delivered
	Seq        int      `json:"seq"`       // Sequence number
}

// Optin
// This callback will occur when the Send-to-Messenger plugin has been tapped.
type Optin struct {
	Sender User
	Ref    string `json:"ref"` // data-ref parameter that was defined with the entry point
}

// Read
// This callback will occur when a message a page has sent has been read by the user.
type Read struct {
	Sender    User
	Watermark float64 `json:"watermark"` // All messages that were sent before this timestamp were read
	Seq       int     `json:"seq"`       // Sequence number
}

// Payment
// TODO: Payment is still in BETA. Implement later.
// Doc: https://developers.facebook.com/docs/messenger-platform/webhook-reference/payment
type Payment struct {
	Sender User
}

// Checkout update
// TODO: Checkout Update is still in BETA. Implement later.
// Document: https://developers.facebook.com/docs/messenger-platform/webhook-reference/checkout-update
type CheckoutUpdate struct {
	Sender User
}
