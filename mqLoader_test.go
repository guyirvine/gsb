package gsb

import (
	"os"
	"testing"
)

func TestLoadMQs(t *testing.T) {
	os.Setenv("GSB_MQ", "inmem://")
	os.Setenv("GSB_ERRORQ", "inmem://")

	h := new(Host)
	h.loadMq()

	Equals(t, "MqInMemory", getTypeName(h.mq))
	Equals(t, "MqInMemory", getTypeName(h.errorQ))

}
