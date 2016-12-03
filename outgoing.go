package fbbot

// TODO: Audio message = Text + audio
// https://developers.facebook.com/docs/messenger-platform/send-api-reference/audio-attachment

// TODO: File message = Text + file
// https://developers.facebook.com/docs/messenger-platform/send-api-reference/audio-attachment

// TODO: Video message = Text + Video
// https://developers.facebook.com/docs/messenger-platform/send-api-reference/video-attachment

// TextMessage contains only text
type TextMessage struct {
	// Text is text content of the message
	// must be UTF-8, 320 character limit
	// required
	Text string
	Noti string
}

func NewTextMessage(text string) *TextMessage {
	var t TextMessage
	t.Text = text
	t.Noti = NotiRegular
	return &t
}

// ButtonMessage contains text and buttons
type ButtonMessage struct {
	// // Text is text content of the message
	// // must be UTF-8, 320 character limit
	// // not required
	// Text string

	// // MainText is text that appears in main body
	// // required
	// MainText string
	Text string
	Noti string

	Buttons []Button
}

func NewButtonMessage() *ButtonMessage {
	var b ButtonMessage
	b.Noti = NotiRegular
	return &b
}

func (m *ButtonMessage) AddWebURLButton(title, URL string) {
	b := NewWebURLButton(title, URL)
	m.Buttons = append(m.Buttons, b)
}

func (m *ButtonMessage) AddPostbackButton(title, payload string) {
	b := NewPostbackButton(title, payload)
	m.Buttons = append(m.Buttons, b)
}

// TODO
// func (m *ButtonMessage) AddButton(b Button) {}
// func (m *ButtonMessage) AddButtons(bs []Button) {}

// Button
type Button struct {
	Type    string `json:"type"` // web_url or postback
	Title   string `json:"title"`
	URL     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

func NewWebURLButton(title, URL string) Button {
	return Button{
		Type:  "web_url",
		Title: title,
		URL:   URL,
	}
}

func NewPostbackButton(title, payload string) Button {
	return Button{
		Type:    "postback",
		Title:   title,
		Payload: payload,
	}
}

// GenericMessage could contain text, image, title, subtitle, description and buttons.
// Can support multiple bubbles per message and display them as a horizontal list.
type GenericMessage struct {
	// Text is text content of the message
	// must be UTF-8, 320 character limit
	// not required
	Text string // TODO: ? nothing to do with this field?
	Noti string

	Bubbles []Bubble
}

func NewGenericMessage() *GenericMessage {
	var g GenericMessage
	g.Noti = NotiRegular
	return &g
}

// Bubble represents ...
type Bubble struct {
	// Title is title of the bubble
	// required
	Title string `json:"title"`

	// SubTitle is buble subtitle
	// not required
	SubTitle string `json:"subtitle"`

	// ItemURL is URL opened when bubble is tapped
	// not required
	ItemURL string `json:"item_url"`

	// ImageURL is URL of bubble image
	// not required
	ImageURL string `json:"image_url"`

	// Buttons are buttons that appear as call-to-actions
	// not required
	Buttons []Button `json:"buttons"`
}

// TODO: Send image from file: https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment

// ImageMessage contains text and image
// Supported formats are jpg, png and gif.
type ImageMessage struct {
	Type string `json:"type"`
	Noti string

	// URL is URL of the image
	// required
	URL string
}

func NewImageMessage() *ImageMessage {
	var i ImageMessage
	i.Type = "image"
	i.Noti = NotiRegular
	return &i
}

type ReceiptMessage struct {
	// Text is text content of the message
	// must be UTF-8, 320 character limit
	// not required
	Text string

	Noti string

	// RecipientName is recipient name
	// required
	RecipientName string

	// OrderNumber is order number
	// must be unique
	// required
	OrderNumber string

	// Currency is currency for the order
	// required
	Currency string

	// PaymentMethod Payment method details. This can be a custom string. Ex: 'Visa 1234'
	// required
	PaymentMethod string

	// Timestamp is timestamp of order
	// not required
	Timestamp string

	// OrderURL is URL of order
	// not required
	OrderURL string

	// Items are items in order
	// required
	Items []Item

	// Shipping address
	// not required
	Address Address

	// Summary is Payment summary
	// required
	Summary Summary

	// Adjustments is Payment adjustments
	// not required
	Adjustments []Adjustment
}

func NewReceiptMessage() *ReceiptMessage {
	var r ReceiptMessage
	r.Noti = NotiRegular
	return &r
}

// Item is item in order
type Item struct {
	// Title is item title
	// required
	Title string

	// Subtile of item
	// not required
	Subtitle string

	// Quantity of item
	// not required
	Quantity float64

	// Price of item
	// not required
	Price float64

	// Currency of item
	// not required
	Currency string

	// ImageURL is image URL of item
	// not required
	ImageURL string
}

// Shipping address
type Address struct {
	// Street1 Street Address, line 1
	Street1 string

	// Street2 Street Address, line 2
	Street2 string

	// City
	City string

	// PostalCode Postal code
	PostalCode string

	// State is state abbrevation
	State string

	// Country is Two-letter country abbreviation
	Country string
}

// Summary is Payment summary
type Summary struct {
	// Subtotal
	// not required
	Subtotal float64

	// ShippingCost is cost of shipping
	// not required
	ShippingCost float64

	// TotalTax is total tax
	// not required
	TotalTax float64

	// TotalCost is total cost
	// required
	TotalCost float64
}

// Adjustment is payment adjustment.
// Allows a way to insert adjusted pricing (e.g., sales).
type Adjustment struct {
	// Name is name of adjustment
	Name string

	// Amount is adjusted amount
	Amout float64
}
