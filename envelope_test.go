package gsb

import (
	"fmt"
	"net/url"
	"testing"
)

func TestFormatErorMessage(t *testing.T) {
	err := fmt.Errorf("TestFormatErorMessage")
	str := formatErrorMessage(err)
	Equals(t, "TestFormatErorMessage\nenvelope_test.go:11", str)
}

func TestNewEnvelope(t *testing.T) {
	mqURL, err := url.Parse("beanstalk://localhost/example")
	Ok(t, err)

	env := createEnvelope("msgName", mqURL.String(), []byte("Payload"), "beanstalk://localhost/reply")

	Equals(t, "msgName", env.MessageName)
}

func TestEnvelope(t *testing.T) {
	mqURL, _ := url.Parse("beanstalk://localhost/example")
	env := createEnvelope("msgName", mqURL.String(), []byte("Payload"), "beanstalk://localhost/reply")

	Equals(t, "Payload", string(env.getMsgPayload()))

	env.addError(fmt.Errorf("TestEnvelope"))

	Equals(t, 1, len(env.Errors))
}
