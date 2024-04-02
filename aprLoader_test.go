package gsb

import (
	"os"
	"testing"
)

func TestLoadAPRs(t *testing.T) {
	os.Setenv("GSB_APR_Store", "inmem://")

	h := new(Host)
	h.loadAPRs()

	os.Setenv("GSB_APR_Store", "")

	Equals(t, 1, len(h.aprList))

}
