package gsb

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestHost(t *testing.T) {
	os.Setenv("GSB_SINGLE_LOOP", "Y")
	os.Setenv("GSB_APR_Store", "inmem://")
	os.Setenv("GSB_MQ", "inmem://")
	// log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)

	host := new(Host)
	host.Init()
	host.LoadHandler(new(DummyHandler))

	m := &DummyMessage{}
	host.Send(m)
	host.Start()
}
