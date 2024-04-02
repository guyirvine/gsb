package gsb

import (
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (h *Host) filterAPRList(rawList []string) map[string]string {
	list := map[string]string{}

	for _, v := range rawList {
		if len(v) < 8 {
			continue
		}
		if v[0:8] != "GSB_APR_" {
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

func (h *Host) getAPRList() map[string]string {
	rawList := os.Environ()
	return h.filterAPRList(rawList)
}

func (h *Host) parseAPRURLString(aprString string) *url.URL {
	aprURL, err := url.Parse(aprString)
	if err != nil {
		log.Fatalf("URL provided for MQ not valid. url: %s", aprString)
	}

	return aprURL
}

func (h *Host) getAPR(aprURL *url.URL) *APRDefinition {
	var apr APRDefinition

	switch aprURL.Scheme {
	case "inmem":
		apr = new(APRInMemory)
	case "postgres":
		apr = new(APRPostgres)
	default:
		log.Fatalf("APR type not supported, scheme: %s\nurl: %s", aprURL.Scheme, aprURL.String())
	}

	err := apr.Init(aprURL)
	if err != nil {
		log.Fatal(err)
	}

	return &apr
}

func (h *Host) loadAPRs() {
	h.aprList = map[string]APRDefinition{}

	aprList := h.getAPRList()
	for name, urlString := range aprList {
		aprURL := h.parseAPRURLString(urlString)
		apr := h.getAPR(aprURL)
		h.aprList[name] = *apr
		log.Infof("APR loaded. %s:%s", name, urlString)
	}
}
