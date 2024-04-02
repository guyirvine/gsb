package gsb

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type Mq interface {
	Init(*url.URL) error
	Connect() error
	Send(*Envelope) error
	GetNextMsg() (*Envelope, error)
	GetMqURL() *url.URL
	Commit() error
}

type MQDefinition struct {
	name  string
	mqURL string
	mq    Mq
}

func validateAndReturnMessageName(msg Message) string {
	// len("Message") = 7

	rawName := getTypeName(msg)

	log.Debug("Message processing. rawName: ", rawName)
	if len(rawName) <= 7 {
		log.Fatalf("Message name too short. Suffix must be 'Message'. name: %s", rawName)
	}
	if rawName[len(rawName)-7:] != "Message" {
		log.Fatalf("Message not named correctly. Suffix must be 'Message'. name: %s", rawName)
	}

	return rawName
}
