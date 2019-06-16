package wit

import (
	"github.com/michlabs/fbbot/nlp"
	"github.com/michlabs/gowit"
	log "github.com/sirupsen/logrus"
)

type Adapter struct {
	client *gowit.Client
}

func init() {
	nlp.Register("wit", &Adapter{})
}

func (a *Adapter) Init(config map[string]string) error {
	a.client = gowit.NewClient(config["token"])
	return nil
}

func (a *Adapter) Detect(msg string) (intent string, entities map[string][]string) {
	meaning, err := a.client.Detect(msg)
	if err != nil {
		log.Error("Failed to detect meaning: ", err)
		return intent, entities
	}

	intent = meaning.Intent()
	entities = make(map[string][]string)
	for key, values := range meaning.Entities {
		if key != "intent" {
			for _, value := range values {
				if str, ok := value["value"].(string); ok {
					entities[key] = append(entities[key], str)
				}
			}
		}
	}

	return intent, entities
}
