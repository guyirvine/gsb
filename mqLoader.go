package gsb

import (
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

func parseMQURLString(mqString string) *url.URL {
	mqURL, err := url.Parse(mqString)
	if err != nil {
		log.Fatalf("URL provided for MQ not valid. url: %s", mqString)
	}

	return mqURL
}

func getMQ(mqString string) Mq {
	mqURL := parseMQURLString(mqString)

	var mq Mq

	switch mqURL.Scheme {
	case "inmem":
		mq = new(MqInMemory)
	case "beanstalk":
		mq = new(mqBeanstalk)
	default:
		errorString := fmt.Sprintf("MQ type not supported, scheme: %s\nurl: %s", mqURL.Scheme, mqString)
		log.Fatal(errorString)
	}

	mq.Init(mqURL)

	log.Info("MQ loaded. ", mqString)

	return mq
}

func (h *Host) loadAuditMqs() {
	h.auditIncomingMq = nil
	h.auditOutgoingMq = nil

	auditIncomingMessages := ""
	auditOutgoingMessages := ""
	if os.Getenv("GSB_MQ_AUDIT") != "" {
		auditIncomingMessages = os.Getenv("GSB_MQ_AUDIT")
		auditOutgoingMessages = os.Getenv("GSB_MQ_AUDIT")
	}
	if os.Getenv("GSB_MQ_AUDIT_INCOMING") != "" {
		auditIncomingMessages = os.Getenv("GSB_MQ_AUDIT_INCOMING")
	}
	if os.Getenv("GSB_MQ_AUDIT_OUTGOING") != "" {
		auditOutgoingMessages = os.Getenv("GSB_MQ_AUDIT_OUTGOING")
	}

	if auditIncomingMessages != "" {
		h.auditIncomingMq = getMQ(auditIncomingMessages)
	}
	if auditOutgoingMessages != "" {
		h.auditOutgoingMq = getMQ(auditOutgoingMessages)
	}
}

func (h *Host) loadMq() {
	h.mq = getMQ(getStrEnv("GSB_MQ", "inmem://"))
	log.Info("Main Q loaded. ")

	h.errorQ = getMQ(getStrEnv("GSB_ERRORQ", "inmem://"))
	log.Info("Error Q loaded. ")

	h.loadAuditMqs()
}
