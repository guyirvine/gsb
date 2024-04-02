package main

import (
	"fmt"
	"gsb"
	"os"

	log "github.com/sirupsen/logrus"
)

type ReplyHandler struct {
	Host *gsb.Host
}

type ReplyMessage struct {
	Label string
}

type ResponseMessage struct {
	Label string
}

func (m *ResponseMessage) GetPayload() string {
	return m.Label
}
func (m *ReplyMessage) GetPayload() string {
	return m.Label
}

func (h *ReplyHandler) GetMessage() gsb.Message {
	return &ReplyMessage{}
}

func (h *ReplyHandler) Init() error {
	return nil
}

func (h *ReplyHandler) Handle(msg gsb.Message) error {
	msg2 := msg.(*ReplyMessage)

	fmt.Printf("ReplyHandler.Handlle. *******   %s\n", msg2.Label)

	h.Host.Reply(&ResponseMessage{"Successful Response"})

	return nil
}

func main() {
	os.Setenv("GSB_APR_Store", "inmem://")
	os.Setenv("GSB_MQ", "beanstalk://localhost/agent_reply")
	replyUrl := "beanstalk://localhost/agent_response"
	// os.Setenv("GSB_MQ", "inmem://")
	log.SetLevel(log.InfoLevel)
	log.SetLevel(log.DebugLevel)

	host := new(gsb.Host)
	host.Init()
	host.LoadHandler(new(ReplyHandler))

	agent := new(gsb.Agent)
	agent.Send(&ReplyMessage{"ttttt"}, os.Getenv("GSB_MQ"), replyUrl)

	host.Start()

	var msg ReplyMessage
	_, err := agent.CheckForReply(replyUrl, &msg)
	if err != nil {
		errorString := fmt.Sprintf("main.CheckForReply.err: %v\n", err)
		log.Error(errorString)
	}

}
