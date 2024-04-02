package gsb

import (
	"os"
	"testing"
)

func TestLoadMQRegistry(t *testing.T) {
	h := new(Host)

	os.Setenv("GSB_MSG_DummyMessage", "inmem://")
	os.Setenv("GSB_MSG_DummyMessage2", "inmem://")
	h.loadMQRegistry()

	Equals(t, 2, len(h.mqList))
	Equals(t, "MqInMemory", getTypeName(h.mqList["DummyMessage"].mq))
}
