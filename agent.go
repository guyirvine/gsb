package gsb

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// A means for a stand-alone process to interact with the bus, without being
// a full rservicebus application

type Agent struct {
}

func (a *Agent) Send(msg Message, mqUrl string, returnAddress string) error {
	os.Setenv(fmt.Sprintf("GSB_MSG_%s", getTypeName(msg)), mqUrl)

	fmt.Printf("Wiggle.1. msg: %s, url: %s\n", getTypeName(msg), mqUrl)

	h := new(Host)
	h.loadMQRegistry()

	env, err := createEnvelopeFromMessage(msg)
	if err != nil {
		log.Debugf("Host.Send.Error createEnvelopeFromMessage msg, %v\nError: %v", msg, err)
		return err
	}
	env.ReplyMqURLString = returnAddress

	h.SendEnvelope(env)

	return nil
}

func (a *Agent) CheckForReply(mqUrl string, msg Message) (Message, error) {
	log.Debugf("Agent.CheckForReply. mqUrl: %s, msg.typeName: %s", mqUrl, getTypeName(msg))

	mq := getMQ(mqUrl)

	env, err := mq.GetNextMsg()
	if err != nil {
		log.Debugf("agent.CheckForReply.GetNextMsg.err: %v\n", err)
		return nil, err
	}

	if env == nil {
		return nil, nil
	}

	log.Debugf("agent.CheckForReply.env: %v\n", env)
	msg, err = getMsg(env.getMsgPayload(), msg)
	if err != nil {
		errorString := fmt.Sprintf("agent.CheckForReply.getMsg.err: %v\n", err)
		log.Error(errorString)
		return nil, err
	}

	mq.Commit()

	return msg, nil
}

func NewAgent() *Agent {
	return new(Agent)
}
