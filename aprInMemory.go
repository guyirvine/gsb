package gsb

import (
	"net/url"

	log "github.com/sirupsen/logrus"
)

type APRInMemory struct {
	list map[string]interface{}
}

func (a *APRInMemory) Reset() error {
	a.list = map[string]interface{}{}
	return nil
}

func (a *APRInMemory) Init(*url.URL) error {
	return a.Reset()
}

func (a *APRInMemory) Set(k string, v interface{}) {
	a.list[k] = v
}

func (a *APRInMemory) Get(k string) interface{} {
	return a.list[k]
}

func (a *APRInMemory) Begin() error {
	log.Debug("APRInMemory.Begin")
	return nil
}
func (a *APRInMemory) Rollback() error {
	log.Debug("APRInMemory.Rollback")
	return nil
}
func (a *APRInMemory) Commit() error {
	log.Debug("APRInMemory.Commit")
	return nil
}
