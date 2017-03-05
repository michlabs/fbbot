package fbbot

import (
	"regexp"
)

type Commander struct {
	mapCommandFunc map[string]func(*Bot, *Message, string)
}

func NewCommander() *Commander {
	commandHandler := Commander{}
	commandHandler.mapCommandFunc = make(map[string]func(*Bot, *Message, string))
	return &commandHandler
}

func (h *Commander) Add(name string, f func(*Bot, *Message, string)) {
	h.mapCommandFunc[name] = f
}

func (h *Commander) HandleEcho(bot *Bot, echoMsg *Message) {
	// Do not handle echo from bot itself
	if echoMsg.AppID > 0 {
		return
	}

	name, param := h.extractCommand(echoMsg.Text)
	f, ok := h.mapCommandFunc[name]
	if ok {
		f(bot, echoMsg, param)
	}
}

func (h *Commander) extractCommand(msg string) (name string, param string) {
	re := regexp.MustCompile("^/([a-z]*)(?: (.*))?")
	matches := re.FindStringSubmatch(msg)
	if len(matches) == 3 {
		name = matches[1]
		param = matches[2]
	}
	return name, param
}
