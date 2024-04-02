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
	return fmt.Errorf("Raise dummy error")
}

func main() {
	os.Setenv("GSB_APR_Store", "inmem://")
	os.Setenv("GSB_MQ", "beanstalk://localhost/main")
	os.Setenv("GSB_ERRORQ", "beanstalk://localhost/error")
	log.SetLevel(log.InfoLevel)
	log.SetLevel(log.DebugLevel)

	host := new(gsb.Host)
	host.Init()
	host.LoadHandler(new(ExampleHandler))

	host.Send(&ExampleMessage{"ttttt"})
	host.Start()
}
