package nlp

import (
	"fmt"
)

var usedAdapter Client
var adapters map[string]Client = make(map[string]Client)

type Client interface {
	Init(config map[string]string) error
	Detect(msg string) (intent string, entities map[string][]string)
}

func Register(name string, client Client) {
	adapters[name] = client
}

func Use(adapter string, config map[string]string) error {
	if _, ok := adapters[adapter]; !ok {
		return fmt.Errorf("Adapter %s is not exist", adapter)
	}
	usedAdapter = adapters[adapter]
	return usedAdapter.Init(config)
}

// Return detected intent and entities
func Detect(msg string) (intent string, entities map[string][]string) {
	return usedAdapter.Detect(msg)
}
