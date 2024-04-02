package gsb

import "net/url"

// APP Resource
type APRDefinition interface {
	Init(*url.URL) error
	Reset() error
	Begin() error
	Rollback() error
	Commit() error
}
