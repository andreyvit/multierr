// Package multierr provides an idiomatic Append function for errors,
// allowing to combine multiple errors into one.
//
// Basically the whole API is:
//
//   err = multierr.Append(err, someFuncThatFails())
//
// Among other use cases, this is especially useful in failable deferred funcs
// to avoid silently ignoring cleanup errors:
//
//   func sprinkleMagicDust() (err error) {
//   	d := &Dummy{}
//   	defer func() {
//   		err = multierr.Append(err, d.Close())
//   	}()
//
//   	return d.DoSomeMagic()
//   }
//
// Unlike many other multierror packages, this one does not even expose its
// multierror type, and will only use it when you actually end up with more than
// a single error to return.
//
// The returned multierror type supports errors.Is and errors.As by delegating
// to each of the suberrors it contains.
package multierr

import (
	"errors"
	"fmt"
	"strings"
)

// multi is the type returned when Append needs to combine multiple errors;
// it will always have at least 2 items.
type multi []error

func (m multi) Error() string {
	switch len(m) {
	case 0:
		panic("multierr.multi does not support zero errors")
	case 1:
		return m[0].Error()
	default:
		return FormatMessage([]error(m))
	}
}

func (m multi) As(target interface{}) bool {
	for _, err := range m {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (m multi) Is(target error) bool {
	for _, err := range m {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// Append joins the given error values into a single error. If both are non-nil,
// wraps them into a multierror type (which is a typed []error slice),
// otherwise returns one of the arguments.
//
//   Append(nil, nil) == nil
//   Append(nil, someErr) == someErr
//   Append(someErr, nil) == someErr
//   Append(someErr, anotherErr) == []error{someErr, anotherErr}
//
// Either (or both) of the arguments can be multierror types, which are properly
// joined together. In this case, Append assumes that it's fine to modify
// either of the arguments.
func Append(dest error, err error) error {
	if err == nil {
		return dest
	} else if dest == nil {
		return err
	} else if md, ok := dest.(multi); ok {
		if ms, ok := err.(multi); ok {
			return append(md, ms...)
		} else {
			return append(md, err)
		}
	} else {
		if ms, ok := err.(multi); ok {
			return append(multi{dest}, ms...)
		} else {
			return multi{dest, err}
		}
	}
}

// ForEach calls f with each suberror in the given error.
// If err is not a multierror type, calls f(err).
// If err is nil, does not call f.
func ForEach(err error, f func(err error)) {
	if err == nil {
		// nop
	} else if m, ok := err.(multi); ok {
		for _, err := range m {
			f(err)
		}
	} else {
		f(err)
	}
}

// Returns the number of suberrors within err.
// If err is not a multierror type, returns 1.
// If err is nil, returns 0.
func Len(err error) int {
	if err == nil {
		return 0
	} else if m, ok := err.(multi); ok {
		return len(m)
	} else {
		return 1
	}
}

// Returns all suberrors within err.
// If err is not a multierror type, returns []error{err}.
// If err is nil, returns nil.
func All(err error) []error {
	var errs []error
	ForEach(err, func(err error) {
		errs = append(errs, err)
	})
	return errs
}

// FormatMessage is a function used to format a string with multiple error messages.
// You can replace it if your project calls for a different format.
// Note that this is a global setting and should be left to the end user to decide.
var FormatMessage func(errs []error) string = DefaultFormatMessage

// DefaultFormatMessage performs the default formatting of multiple error messages.
func DefaultFormatMessage(errs []error) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d errors occurred:\n", len(errs))
	for i, err := range errs {
		if i > 0 {
			buf.WriteByte('\n')
		}
		s := err.Error()
		fmt.Fprintf(&buf, "(%d) %s", i+1, strings.ReplaceAll(s, "\n", "\n\t"))
	}
	return buf.String()
}
