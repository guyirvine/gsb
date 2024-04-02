package gsb

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	log "github.com/sirupsen/logrus"
)

type mqBeanstalk struct {
	b              *beanstalk.Conn
	tube           *beanstalk.Tube
	tubeSet        *beanstalk.TubeSet
	queueName      string
	sendTimeout    time.Duration
	reserveTimeout time.Duration
	msgId          uint64
	mqURL          *url.URL
}

func (mq *mqBeanstalk) GetMqURL() *url.URL {
	return mq.mqURL
}

func (mq *mqBeanstalk) getTimeoutFromString(inputString string, defaultValue time.Duration) time.Duration {
	if inputString == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(inputString, 10, 64)

	if err != nil {
		return defaultValue
	}

	return time.Duration(intValue) * time.Second
}

func (mq *mqBeanstalk) Connect() error {
	var err error

	mq.b = nil
	mq.tube = nil
	mq.tubeSet = nil
	mq.msgId = 0

	hostName := mq.mqURL.Hostname()
	if hostName == "" {
		hostName = "127.0.0.1"
	}
	port := mq.mqURL.Port()
	if port == "" {
		port = "11300"
	}

	address := fmt.Sprintf("%s:%s", hostName, port)
	mq.b, err = beanstalk.Dial("tcp", address)

	if err != nil {
		log.Errorf("MQ. Error dialing beanstalk. address: %s\nError from beanstalk: %s", address, err.Error())
		return err
	}

	queueName := mq.mqURL.Path
	if queueName == "/" || queueName == "" {
		queueName = "main"
	} else {
		queueName = queueName[1:]
	}
	mq.queueName = queueName
	mq.tube = beanstalk.NewTube(mq.b, mq.queueName)
	mq.tubeSet = beanstalk.NewTubeSet(mq.b, mq.queueName)

	query := mq.mqURL.Query()
	mq.sendTimeout = mq.getTimeoutFromString(query.Get("sendTimeout"), 120)
	mq.reserveTimeout = mq.getTimeoutFromString(query.Get("sendTimeout"), 1)

	log.Debugf("mqBeanstalk.loaded. Queuename: %s, url: %s", mq.queueName, mq.mqURL.String())

	return nil
}

func (mq *mqBeanstalk) Init(mqURL *url.URL) error {
	mq.mqURL = mqURL

	return mq.Connect()
}

func (mq *mqBeanstalk) Send(env *Envelope) error {
	payload, err := json.Marshal(env)
	if err != nil {
		return err
	}

	log.Debug("MQ.Beanstalk.Send. env: ", env)

	_, err = mq.tube.Put([]byte(payload), 1, 0, mq.sendTimeout)

	return err
}

func (mq *mqBeanstalk) GetNextMsg() (*Envelope, error) {
	var err error
	var payload []byte

	mq.msgId, payload, err = mq.tubeSet.Reserve(mq.reserveTimeout)

	if err != nil {
		return nil, err
	}

	var env Envelope
	if err = json.Unmarshal(payload, &env); err != nil {
		return nil, err
	}

	log.Debug("MQ.mqBeanstalk.env:", env)
	log.Debug("MQ.mqBeanstalk.env.MsgPayload:", string(env.getMsgPayload()))

	return &env, nil
}

func (mq *mqBeanstalk) Commit() error {
	return mq.b.Delete(mq.msgId)
}
