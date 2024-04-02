package gsb

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func filterMQList(rawList []string) map[string]string {
	list := map[string]string{}

	for _, v := range rawList {
		if len(v) < 8 {
			continue
		}
		if v[0:8] != "GSB_MSG_" {
			continue
		}
		kvs := strings.Split(v, "=")
		if len(kvs) != 2 {
			log.Fatalf("App Resource Environment Variable not formatted correclty. env: %s", v)
		}
		keyName := kvs[0][8:]
		list[keyName] = kvs[1]
	}

	return list
}

func (h *Host) loadMQRegistry() {
	h.mqList = map[string]MQDefinition{}

	rawList := os.Environ()
	mqList := filterMQList(rawList)
	log.Debugf("mqRegistryLoader.filterMQList.mqList: %v", mqList)

	for name, urlString := range mqList {
		mq := getMQ(urlString)
		mqd := &MQDefinition{name, urlString, mq}
		h.mqList[name] = *mqd
		log.Infof("Mq loaded. %s:%s", name, urlString)
	}
}
