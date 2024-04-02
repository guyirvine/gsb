package gsb

import (
	"testing"
)

func TestLoadHandlers(t *testing.T) {
	h := new(Host)
	h.loadHandlers()

	Equals(t, 1, len(h.handlers))
}
