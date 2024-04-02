package main

import (
	"fmt"
	"gsb"
	"os"

	log "github.com/sirupsen/logrus"
)

type RegistryHandler struct {
	Host *gsb.Host
}

type RegistryMessage struct {
	Label string
}

type RegistrySendMessage struct {
	Label string
}

func (m *RegistryMessage) GetPayload() string {
	return m.Label
}
func (m *RegistrySendMessage) GetPayload() string {
	return m.Label
}

func (h *RegistryHandler) GetMessage() gsb.Message {
	return &RegistryMessage{}
}

func (h *RegistryHandler) Init() error {
	return nil
}

func (h *RegistryHandler) Handle(msg gsb.Message) error {
	msg2 := msg.(*RegistryMessage)

	fmt.Printf("RegistryHandler.Handlle. *******   %s\n", msg2.Label)

	h.Host.Send(&RegistrySendMessage{"Successful Response"})

	return nil
}

func main() {
	os.Setenv("GSB_APR_Store", "inmem://")
	os.Setenv("GSB_MQ", "beanstalk://localhost/registry")
	os.Setenv("GSB_MSG_RegistrySendMessage", "beanstalk://localhost/registrysend")
	log.SetLevel(log.InfoLevel)
	//	log.SetLevel(log.DebugLevel)

	host := new(gsb.Host)
	host.Init()
	host.LoadHandler(new(RegistryHandler))

	agent := new(gsb.Agent)
	agent.Send(&RegistryMessage{"ttttt"}, os.Getenv("GSB_MQ"), "")

	host.Start()

	var msg RegistrySendMessage
	_, err := agent.CheckForReply("beanstalk://localhost/registrysend", &msg)
	if err != nil {
		errorString := fmt.Sprintf("main.CheckForRegistrySend.err: %v\n", err)
		log.Fatal(errorString)
	}

	log.Debugf("main.RegistrySendMessage. msg: %v\n", msg)
	log.Printf("main.RegistrySendMessage. payload: %s\n", msg.GetPayload())
}
