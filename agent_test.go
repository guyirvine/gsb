package gsb

import (
	"os"
	"testing"
)

func TestSend(t *testing.T) {
	os.Setenv("GSB_APR_Store", "inmem://")
	os.Setenv("GSB_MQ", "beanstalk://localhost/agent_reply")
	replyUrl := "beanstalk://localhost/agent_response"

	host := new(Host)
	host.Init()
	host.LoadHandler(new(DummyHandler))

	NewAgent().Send(&DummyMessage{}, os.Getenv("GSB_MQ"), replyUrl)

	//	Equals(t, "mqInMemory", getTypeName(h.mq))
	//	Equals(t, "mqInMemory", getTypeName(h.errorQ))

}
