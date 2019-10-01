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
	var taggedErrors map[string][]string

	for _, err := range e.errs {
		if ke, ok := err.(taggedError); ok {
			if taggedErrors == nil {
				taggedErrors = make(map[string][]string)
			}
			taggedErrors[ke.Error()] = append(taggedErrors[ke.Error()], ke.key)
		} else {
			lines = append(lines, err.Error())
		}
	}

	var orderedKeyedErrors []string
	for err := range taggedErrors {
		orderedKeyedErrors = append(orderedKeyedErrors, err)
	}
	sort.Strings(orderedKeyedErrors)
	for _, err := range orderedKeyedErrors {
		lines = append(lines, fmt.Sprintf("%s (%v)", err, strings.Join(taggedErrors[err], ", ")))
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

// Unfold returns the underlying list of errors wrapped in a multierror.
// If err is not a multierror, then a singleton list is returned.
func Unfold(err error) []error {
	if me, ok := err.(*Error); ok {
		return me.errs
	} else {
		return []error{err}
	}
}

// Uniq deduplicates a list of errors
func Uniq(errs []error) []error {
	type groupingKey struct {
		msg    string
		tagged bool
	}
	var ordered []groupingKey
	grouped := map[groupingKey][]error{}

	for _, err := range errs {
		msg, tag := TaggedError(err)
		key := groupingKey{
			msg:    msg,
			tagged: tag != "",
		}
		if _, ok := grouped[key]; !ok {
			ordered = append(ordered, key)
		}
		grouped[key] = append(grouped[key], err)
	}

	var res []error
	for _, key := range ordered {
		group := grouped[key]
		err := group[0]
		if key.tagged {
			var tags []string
			for _, e := range group {
				_, tag := TaggedError(e)
				tags = append(tags, tag)
			}
			err = fmt.Errorf("%w (%s)", err, strings.Join(tags, ", "))
		} else {
			if n := len(group); n > 1 {
				err = fmt.Errorf("%w repeated %d times", err, n)
			}
		}
		res = append(res, err)
	}

	return res
}

type TaggableError interface {
	// TaggedError is like Error() but splits the error from the tag.
	TaggedError() (string, string)
}

// TaggedError is like Error() but if err implements TaggedError, it will
// invoke TaggeddError() and return error message and the tag. Otherwise the tag will be empty.
func TaggedError(err error) (string, string) {
	if te, ok := err.(TaggableError); ok {
		return te.TaggedError()
	}
	return err.Error(), ""
}

type taggedError struct {
	error
	key string
}

// Tagged wraps an error with a tag. All errors sharing the same error msg will be grouped together in one entry
// of the multierror along with the list of tags.
func Tagged(key string, err error) error {
	return taggedError{error: err, key: key}
}

func (k taggedError) Unwrap() error {
	return k.error
}

func (k taggedError) TaggedError() (string, string) {
	return k.Error(), k.key
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
