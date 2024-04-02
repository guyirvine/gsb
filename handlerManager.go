package gsb

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// getNameForHandler ensures and strips the Message suffix to line up with the Handler
func (h *Host) getNameForHandler(msgName string) (string, error) {
	// len("Message") = 7
	rawName := msgName

	log.Debug("handlerManager.getNameForHandler. rawName: ", rawName)
	if len(rawName) <= 7 {
		return "", fmt.Errorf("Message name too short. Suffix must be 'Message'. name: %s", rawName)
	}
	if rawName[len(rawName)-7:] != "Message" {
		return "", fmt.Errorf("Message not named correctly. Suffix must be 'Message'. name: %s", rawName)
	}
	name := rawName[0 : len(rawName)-7]

	return name, nil
}

func (h *Host) getHandlerDefinition(env *Envelope) (*HandlerDefinition, error) {
	msgName, err := h.getNameForHandler(env.MessageName)
	if err != nil {
		return nil, err
	}
	log.Debugf("handlerManager.getHanlder.msgName: %s", msgName)
	log.Infof("handlerManager.getHanlder.handlers: %v", h.handlers)

	if h.handlers[msgName] == nil {
		return nil, fmt.Errorf("Host.processEnvelope. Could not find Handler for Message: %s", msgName)
	}

	return h.handlers[msgName], nil
}
