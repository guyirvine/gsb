package gsb

import (
	"testing"
)

func TestHandler(t *testing.T) {
	handler := new(DummyHandler)
	handler.Handle(handler.GetMessage())
}
