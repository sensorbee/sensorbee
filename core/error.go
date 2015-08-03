package core

import (
	"fmt"
	"strings"
)

var (
	errPipeClosed = fmt.Errorf("the pipe is already closed")
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

// Fatal returns true when at least one of errors is a fatal error.
func (b *bulkErrors) Fatal() bool {
	for _, e := range b.errs {
		if IsFatalError(e) {
			return true
		}
	}
	return false
}

// Temporary returns true if all the errors bulkErrors has are temporary errors.
func (b *bulkErrors) Temporary() bool {
	for _, e := range b.errs {
		if !IsTemporaryError(e) {
			return false
		}
	}
	return true
}

// IsFatalError returns true when the given error is fatal. If the error
// implements the following interface, IsFatalError returns the return value of
// Fatal method:
//
//	interface {
//		Fatal() bool
//	}
//
// Otherwise, the error is considered non-fatal and it returns false.
//
// Each component in core package behaves differently when it receives a
// fatal error. All functions or methods have a documentation about how they
// behave on this error, so please read them for details.
func IsFatalError(err error) bool {
	type fatal interface {
		Fatal() bool
	}

	if e, ok := err.(fatal); ok {
		return e.Fatal()
	}
	return false
}

type fatalError struct {
	err error
}

func (f *fatalError) Error() string {
	return f.err.Error()
}

func (f *fatalError) Fatal() bool {
	return true
}

// FatalError decorates the given error so that IsFatalError(FatalError(err))
// returns true even if err doesn't have Fatal method. It will panic if err is
// nil.
func FatalError(err error) error {
	if err == nil {
		panic(fmt.Errorf("the error cannot be nil"))
	}
	return &fatalError{err: err}
}

// IsTemporaryError returns true when the given error is temporary. If the error
// implements the following interface, IsTemporaryError returns the return
// value of Temporary method:
//
//	interface {
//		Temporary() bool
//	}
//
// Otherwise, the error is considered permanent and it returns false.
//
// Each component in core package behaves differently when it receives a
// temporary error. All functions or methods have a documentation about how they
// behave on this error, so please read them for details.
func IsTemporaryError(err error) bool {
	type temporary interface {
		Temporary() bool
	}

	if e, ok := err.(temporary); ok {
		return e.Temporary()
	}
	return false
}

type temporaryError struct {
	err error
}

func (t *temporaryError) Error() string {
	return t.err.Error()
}

func (t *temporaryError) Temporary() bool {
	return true
}

// TemporaryError decorates the given error so that
// IsTemporaryError(TemporaryError(err)) returns true even if err doesn't have
// Temporary method. It will panic if err is nil.
func TemporaryError(err error) error {
	if err == nil {
		panic(fmt.Errorf("the error cannot be nil"))
	}
	return &temporaryError{err: err}
}

// TODO: add a hybrid error interface having all possible methods which can
// customize behavior by setting flags.
