package gsb

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Host struct {
	handlers        map[string]*HandlerDefinition
	mq              Mq
	errorQ          Mq
	auditIncomingMq Mq
	auditOutgoingMq Mq
	aprList         map[string]APRDefinition
	mqList          map[string]MQDefinition
	MaxRetries      int
	env             *Envelope
}

func getMsg(msgPayload []byte, msg Message) (Message, error) {
	err := json.Unmarshal(msgPayload, &msg)
	log.Debug("envelope.getMsg. ", msg)

	return msg, err
}

func (h *Host) processMessage(msg Message, hd *HandlerDefinition) error {
	// Begin Transaction for all App Resources
	for _, apr := range hd.aprList {
		if err := apr.Begin(); err != nil {
			log.Debugf("Error raised issueing apr.Begin(), %s, for msg %v.\nError: %v", getTypeName(apr), msg, err)
			return err
		}
	}

	err := hd.handler.Handle(msg)
	if err != nil {
		log.Debugf("Error raised processing Handler, %s, with msg, %v.\nError: %v", hd.name, msg, err)
		// Rollback Transaction for all App Resources
		for _, apr := range hd.aprList {
			apr.Rollback()
		}

		return err
	}

	// Commit Transaction for all App Resources
	for _, apr := range hd.aprList {
		apr.Commit()
	}
	if h.auditIncomingMq != nil {
		h.auditIncomingMq.Send(h.env)
	}
	return nil
}

func (h *Host) processEnvelope() {
	hd, err := h.getHandlerDefinition(h.env)
	if err != nil {
		h.env.addError(err)
		h.errorQ.Send(h.env)
		return
	}
	log.Debugf("host.processEnvelope.Handler, %s. APR Count: %d", hd.name, len(hd.aprList))

	msg, err := getMsg(h.env.getMsgPayload(), hd.handler.GetMessage())
	if err != nil {
		log.Debugf("Host.processEnvelope.getMsg.err: %v\n", err)
		h.env.addError(err)
		h.errorQ.Send(h.env)
	}

	success := false
	for retry := 0; retry < h.MaxRetries; retry++ {
		log.Debugf("Host.processEnvelope. Processing message, %s, retry: %d\n", hd.name, retry)
		err = h.processMessage(msg, hd)
		log.Debugf("Host.processEnvelope. In loop, Message: %s, err: %v\n", hd.name, err)
		if err == nil {
			success = true
			break
		}

		// Something went wrong so reset all APR's before trying again
		for _, apr := range hd.aprList {
			apr.Reset()
		}
	}
	log.Debugf("Host.processEnvelope. Message: %s, success: %v, err: %v\n", hd.name, success, err)
	if !success {
		log.Debugf("Host.processEnvelope. Processing error. Message: %s, err: %v\n", hd.name, err)
		h.env.addError(err)
		h.errorQ.Send(h.env)
	}
}

func (h *Host) MainLoop() {
	singleLoop := os.Getenv("GSB_SINGLE_LOOP") == "Y"
	shouldKeepLooping := !singleLoop
	loop := true
	maxRetry := 3
	retry := maxRetry

	for loop {
		env, err := h.mq.GetNextMsg()
		if err != nil {

			// reset the connection to the message queue
			h.mq.Connect()

			if retry = retry - 1; retry == 0 {
				log.Fatal(err)
			}
		}

		retry = maxRetry
		if env == nil {
			loop = shouldKeepLooping
			continue
		}

		h.env = env
		h.processEnvelope()
		err = h.mq.Commit()
		if err != nil {
			log.Errorf("Error occured commiting Envelope: %v\nError: %v", env, err)
		}

		log.Debug("Main Loop")
	}
}

func (h *Host) Init() {
	h.loadAPRs()
	h.loadMq()
	h.loadHandlers()
	h.loadMQRegistry()
	h.MaxRetries = getIntEnv("GSB_MAX_RETRIES", 3)
}

func (h *Host) Start() {
	h.MainLoop()
}

func (h *Host) SendEnvelope(env *Envelope) error {
	log.Debugf("Host.SendEnvelope.MessageName: %s", env.MessageName)
	mq := h.mqList[env.MessageName].mq
	if mq == nil {
		log.Debugf("Host.Send. Mq not mapped for Message. registry: %v", h.mqList)
		mq = h.mq
	}

	env.MqURLString = mq.GetMqURL().String()

	mq.Send(env)

	return nil
}

func (h *Host) Send(msg Message) error {
	log.Debug("host.Send.m.typeName: ", getTypeName(msg))

	env, err := createEnvelopeFromMessage(msg)
	if err != nil {
		log.Debugf("Host.Send.Error createEnvelopeFromMessage msg, %v\nError: %v", msg, err)
		return err
	}

	h.SendEnvelope(env)

	return nil
}

func (h *Host) Reply(msg Message) error {
	log.Debugf("Host.Reply.msg.typeName: %s, url: %s", getTypeName(msg), h.env.ReplyMqURLString)

	mq := getMQ(h.env.ReplyMqURLString)

	msgName := validateAndReturnMessageName(msg)
	msgPayload, err := json.Marshal(msg)
	if err != nil {
		log.Debugf("Host.Reply.Error marshalling msg, %v\nError: %v", msg, err)
		return err
	}
	env := createEnvelope(msgName, mq.GetMqURL().String(), msgPayload, "")

	mq.Send(env)

	return nil
}
