package main

import (
	"fmt"
	"gsb"
	"os"

	log "github.com/sirupsen/logrus"
)

type ExampleHandler struct {
	Store *gsb.APRInMemory
}

type ExampleMessage struct {
	Label string
}

func (m *ExampleMessage) GetPayload() string {
	return m.Label
}

func (h *ExampleHandler) GetMessage() gsb.Message {
	return &ExampleMessage{}
}

func (h *ExampleHandler) Init() error {
	return nil
}

func (h *ExampleHandler) Handle(msg gsb.Message) error {
	msg2 := msg.(*ExampleMessage)

	h.Store.Set("name", "12345")
	val := h.Store.Get("name")
	fmt.Printf("*******   ExampleHandler.val: %v\n", val)

	fmt.Printf("*******   %s\n", msg2.Label)

	return nil
}

func main() {
	os.Setenv("GSB_APR_Store", "inmem://")
	os.Setenv("GSB_MQ", "beanstalk://")
	// os.Setenv("GSB_MQ", "inmem://")
	// log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)

	host := new(gsb.Host)
	host.Init()
	host.LoadHandler(new(ExampleHandler))

	m := &ExampleMessage{"rrrrr"}
	host.Send(m)
	host.Send(&ExampleMessage{"ttttt"})
	host.Start()

}
