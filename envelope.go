package gsb

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Message interface {
	GetPayload() string
}

type Error struct {
	ErrorMsg  string
	Timestamp time.Time
}

type Envelope struct {
	MsgId            string
	MsgString        string
	Errors           []Error
	MessageName      string
	ReceivedAt       time.Time
	MqURLString      string
	ReplyMqURLString string
}

// formatErrorMessage formats the error message adding to the envelope
func formatErrorMessage(error error) string {
	if error == nil {
		return ""
	}

	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", filepath.Base(file), line)

	formattedError := fmt.Sprintf("%s\n%s", error.Error(), location)

	return formattedError
}

func (env *Envelope) addError(err error) {
	log.Debugf("envelope.addError. Error, %v\n", err)
	e := Error{formatErrorMessage(err), time.Now()}
	env.Errors = append(env.Errors, e)
}

func (env *Envelope) getMsgPayload() []byte {
	return []byte(env.MsgString)
}

func createEnvelope(msgName string, mqURLString string, payload []byte, replyMqURLString string) *Envelope {
	env := new(Envelope)
	env.MqURLString = mqURLString
	env.MsgId = uuid.NewString()
	env.Errors = []Error{}
	env.MessageName = msgName
	env.MsgString = string(payload)
	env.ReceivedAt = time.Now()
	env.ReplyMqURLString = replyMqURLString

	return env
}

func createEnvelopeFromMessage(msg Message) (*Envelope, error) {
	msgName := validateAndReturnMessageName(msg)
	msgPayload, err := json.Marshal(msg)
	if err != nil {
		log.Debugf("Host.Send.Error marshalling msg, %v\nError: %v", msg, err)
		return nil, err
	}

	env := createEnvelope(msgName, "", msgPayload, "")
	return env, nil
}
