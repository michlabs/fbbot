package fptai

import (
	log "github.com/Sirupsen/logrus"
	fptai "github.com/fpt-corp/fptai-sdk-go"
	"github.com/michlabs/fbbot/nlp"
)

type FPTAI struct {
	app *fptai.Application
}

func init() {
	nlp.Register("fptai", &FPTAI{})
}

// Example of config
// config := map[string]string{
// 	"username": "your_username",
// 	"password": "your_password",
// 	"application_code": "your_application_code",
// }
func (fpt *FPTAI) Init(config map[string]string) error {
	client, err := fptai.NewClient(config["username"], config["password"])
	if err != nil {
		return err
	}
	fpt.app = client.GetApp(config["application_code"])

	return nil
}

func (fpt *FPTAI) Detect(msg string) (intent string, entities map[string][]string) {
	ir, err := fpt.app.Recognize(msg)
	if err != nil {
		log.Errorf("failed to detect message: %s. Error: %s", msg, err.Error())
		return intent, entities
	}
	intent = ir.Intent

	return intent, entities
}