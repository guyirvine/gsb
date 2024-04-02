package gsb

import log "github.com/sirupsen/logrus"

type Handler interface {
	GetMessage() Message
	Handle(Message) error
	Init() error
}

type HandlerDefinition struct {
	name    string
	handler Handler
	aprList []APRDefinition
}

type DummyHandler struct {
}
type DummyMessage struct {
}

func (m *DummyMessage) GetPayload() string {
	return "Forty"
}

func (h *DummyHandler) GetMessage() Message {
	return &DummyMessage{}
}

func (h *DummyHandler) Init() error {
	return nil
}

func (h *DummyHandler) Handle(msg Message) error {
	payload := msg.GetPayload()
	log.Info(payload)

	return nil
}
