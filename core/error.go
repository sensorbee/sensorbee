package core

import (
	"strings"
)

type bulkErrors struct {
	errs []error
}

func (b *bulkErrors) append(err error) {
	b.errs = append(b.errs, err)
}

func (b *bulkErrors) returnError() error {
	switch len(b.errs) {
	case 0:
		return nil
	case 1:
		return b.errs[0]
	default:
		return b
	}
}

func (b *bulkErrors) Error() string {
	// TODO: format could be improved
	res := make([]string, len(b.errs))
	for i, e := range b.errs {
		res[i] = e.Error()
	}
	return `"` + strings.Join(res, `", "`) + `"`
}
