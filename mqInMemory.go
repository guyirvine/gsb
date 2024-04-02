package gsb

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type MqInMemory struct {
	list  []*Envelope
	mqURL *url.URL
}

func (mq *MqInMemory) GetMqURL() *url.URL {
	return mq.mqURL
}

func (mq *MqInMemory) Connect() error {
	mq.list = []*Envelope{}

	return nil
}

func (mq *MqInMemory) Init(mqURL *url.URL) error {
	mq.mqURL = mqURL

	mq.Connect()

	return nil
}

func (mq *MqInMemory) Send(env *Envelope) error {
	mq.list = append(mq.list, env)
	return nil
}

func (mq *MqInMemory) GetNextMsg() (*Envelope, error) {
	if len(mq.list) < 1 {
		return nil, nil
	}

	env := mq.list[0]
	log.Debugf("mqInMemory.GetNextMsg.env: %v\n", env)
	log.Debugf("mqInMemory.GetNextMsg.MsgPayload: %v\n", string(env.getMsgPayload()))

	return mq.list[0], nil
}

func (mq *MqInMemory) Commit() error {
	if len(mq.list) < 1 {
		return nil
	}

	mq.list = mq.list[1:]
	return nil
}
