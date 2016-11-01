package main 

import (
	"github.com/michlabs/fbbot"
)

const PORT int = 8080
const VERIFYTOKEN string = "your_secure_token"
const PAGEACCESSTOKEN string = "your_beloved_page_access_token"

func main() {
	bot := fbbot.New(PORT, VERIFYTOKEN, PAGEACCESSTOKEN)
	bot.HandleMessage(MessageHandler)
	bot.Run()
}

func MessageHandler(bot *fbbot.Bot, msg *fbbot.Message) {
	// Echo... echo...
	m := fbbot.NewTextMessage(msg.Text)
	bot.Send(msg.Sender, m)
}