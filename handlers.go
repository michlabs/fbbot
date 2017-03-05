package fbbot

type MessageHandler interface {
	HandleMessage(*Bot, *Message)
}

type PostbackHandler interface {
	HandlePostback(*Bot, *Postback)
}

type DeliveryHandler interface {
	HandleDelivery(*Bot, *Delivery)
}

type OptinHandler interface {
	HandleOptin(*Bot, *Optin)
}

type ReadHandler interface {
	HandleRead(*Bot, *Read)
}

type EchoHandler interface {
	HandleEcho(*Bot, *Message)
}

type CheckoutUpdateHandler interface {
	HandleCheckoutUpdate(*Bot, *CheckoutUpdate)
}

type PaymentHandler interface {
	HandlePayment(*Bot, *Payment)
}
