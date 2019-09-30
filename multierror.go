package multierror

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Error bundles multiple errors and make them obey the error interface
type Error struct {
	errs []error
}

func (e *Error) Error() string {
	buf := bytes.NewBuffer(nil)

	var keyedErrors map[string][]string

	fmt.Fprintf(buf, "%d errors occurred:", len(e.errs))
	for _, err := range e.errs {
		if ke, ok := err.(keyedError); ok {
			if keyedErrors == nil {
				keyedErrors = make(map[string][]string)
			}
			keyedErrors[ke.Error()] = append(keyedErrors[ke.Error()], ke.key)
		} else {
			fmt.Fprintf(buf, "\n%v", err)
		}
	}

	var orderedKeyedErrors []string
	for err := range keyedErrors {
		orderedKeyedErrors = append(orderedKeyedErrors, err)
	}
	sort.Strings(orderedKeyedErrors)
	for _, err := range orderedKeyedErrors {
		fmt.Fprintf(buf, "\n%s (%v)", err, strings.Join(keyedErrors[err], ", "))
	}

	return buf.String()
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
		return &Error{errs}
	}
	switch err := err.(type) {
	case *Error:
		err.errs = append(err.errs, errs...)
		return err
	default:
		return &Error{append([]error{err}, errs...)}
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
