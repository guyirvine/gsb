package gsb

import (
	"testing"
)

func TestHandlerManageExists(t *testing.T) {
	h := new(Host)
	h.loadHandlers()

	env := new(Envelope)
	env.MessageName = "DummyMessage"
	hd, err := h.getHandlerDefinition(env)

	Equals(t, nil, err)
	Equals(t, "Dummy", hd.name)
}

func TestHandlerManageNotExists(t *testing.T) {
	h := new(Host)
	h.loadHandlers()

	env := new(Envelope)
	env.MessageName = "Dummy2Message"
	_, err := h.getHandlerDefinition(env)

	Equals(t, true, err != nil)
}
