package gsb

import (
	"net/url"
	"testing"
)

func TestMqInMemory(t *testing.T) {
	mq := new(MqInMemory)

	mqURL, _ := url.Parse("inmem://")

	mq.Init(mqURL)

	env, err := mq.GetNextMsg()
	Equals(t, (*Envelope)(nil), env)
	Equals(t, nil, err)

	env = new(Envelope)
	env.MessageName = "Test"
	mq.Send(env)

	env2, err := mq.GetNextMsg()
	Equals(t, env2.MessageName, env.MessageName)
	Equals(t, nil, err)
}
