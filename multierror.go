package multierror

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Error bundles multiple errors and make them obey the error interface
type Error struct {
	errs      []error
	formatter Formatter
}

// Formatter allows to customize the rendering of the multierror.
type Formatter func(errs []string) string

var DefaultFormatter = func(errs []string) string {
	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, "%d errors occurred:", len(errs))
	for _, line := range errs {
		fmt.Fprintf(buf, "\n%s", line)
	}

	return buf.String()
}

func (e *Error) Error() string {
	var f Formatter = DefaultFormatter
	if e.formatter != nil {
		f = e.formatter
	}
	var lines []string
	var keyedErrors map[string][]string

	for _, err := range e.errs {
		if ke, ok := err.(keyedError); ok {
			if keyedErrors == nil {
				keyedErrors = make(map[string][]string)
			}
			keyedErrors[ke.Error()] = append(keyedErrors[ke.Error()], ke.key)
		} else {
			lines = append(lines, err.Error())
		}
	}

	var orderedKeyedErrors []string
	for err := range keyedErrors {
		orderedKeyedErrors = append(orderedKeyedErrors, err)
	}
	sort.Strings(orderedKeyedErrors)
	for _, err := range orderedKeyedErrors {
		lines = append(lines, fmt.Sprintf("%s (%v)", err, strings.Join(keyedErrors[err], ", ")))
	}

	return f(lines)
}

// Append creates a new mutlierror.Error structure or appends the arguments to an existing multierror
// err can be nil, or can be a non-multierror error.
//
// If err is nil and errs has only one element, that element is returned.
// I.e. a singleton error is never treated and (thus rendered) as a multierror.
// This also also effectively allows users to just pipe through the error value of a function call,
// without having to first check whether the error is non-nil.
func Append(err error, errs ...error) error {
	if err == nil && len(errs) == 1 {
		return errs[0]
	}
	if len(errs) == 1 && errs[0] == nil {
		return err
	}
	if err == nil {
		return &Error{errs: errs}
	}
	switch err := err.(type) {
	case *Error:
		err.errs = append(err.errs, errs...)
		return err
	default:
		return &Error{errs: append([]error{err}, errs...)}
	}
}

type keyedError struct {
	error
	key string
}

// Keyed wraps an error with a key. All errors sharing the same key will be grouped together in one entry
// of the multierror along with the list of keys.
func Keyed(key string, err error) error {
	return keyedError{error: err, key: key}
}

// WithFormatter sets a custom formatter if err is a multierror.
func WithFormatter(err error, f Formatter) error {
	if me, ok := err.(*Error); ok {
		cpy := *me
		cpy.formatter = f
		return &cpy
	} else {
		return err
	}
}
