package multierror

import (
	"testing"

	"github.com/cesanta/errors"
)

func TestAppend(t *testing.T) {
	var err error
	err = Append(err, errors.Errorf("an error"))
	if err == nil {
		t.Fatal(err)
	}

	if got, want := err.Error(), `an error`; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	err = Append(err, errors.Errorf("another error"))
	if got, want := err.Error(), `2 errors occurred:
an error
another error`; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	err = errors.Errorf("old error")
	err = Append(err, errors.Errorf("new error"))
	if err == nil {
		t.Fatal(err)
	}

	if got, want := err.Error(), `2 errors occurred:
old error
new error`; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestAppendNil(t *testing.T) {
	var err error
	err = Append(err, nil)
	if err != nil {
		t.Errorf("should be nil")
	}
}

func TestAppendNilOnSomething(t *testing.T) {
	err1 := errors.Errorf("test")
	errs := err1
	errs = Append(errs, nil)

	if got, want := errs, err1; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestAppendMultiple(t *testing.T) {
	err1 := errors.Errorf("test")
	var errs error
	errs = Append(nil, err1)
	errs = Append(errs, nil)

	if got, want := errs, err1; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
