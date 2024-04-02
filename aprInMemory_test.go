package gsb

import (
	"net/url"
	"testing"
)

func TestAprInMemory(t *testing.T) {
	apr := new(APRInMemory)

	mqURL, _ := url.Parse("inmem://")

	apr.Init(mqURL)
	apr.Set("key", "value")
	Equals(t, "value", apr.Get("key"))
	apr.Reset()
	Equals(t, nil, apr.Get("key"))
	apr.Begin()
	apr.Commit()
	apr.Rollback()
}
