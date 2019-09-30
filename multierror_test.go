package multierror_test

import (
	"fmt"
	"testing"

	"github.com/mkmik/multierror"
)

func TestAppend(t *testing.T) {
	var err error
	err = multierror.Append(err, fmt.Errorf("an error"))
	if err == nil {
		t.Fatal(err)
	}

	if got, want := err.Error(), `an error`; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	err = multierror.Append(err, fmt.Errorf("another error"))
	if got, want := err.Error(), `2 errors occurred:
an error
another error`; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	err = fmt.Errorf("old error")
	err = multierror.Append(err, fmt.Errorf("new error"))
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
	err = multierror.Append(err, nil)
	if err != nil {
		t.Errorf("should be nil")
	}
}

func TestAppendNilOnSomething(t *testing.T) {
	err1 := fmt.Errorf("test")
	errs := err1
	errs = multierror.Append(errs, nil)

	if got, want := errs, err1; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestAppendMultiple(t *testing.T) {
	err1 := fmt.Errorf("test")
	var errs error
	errs = multierror.Append(nil, err1)
	errs = multierror.Append(errs, nil)

	if got, want := errs, err1; got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func ExampleAppend() {
	var errs error

	errs = multierror.Append(errs, fmt.Errorf("foo"))
	errs = multierror.Append(errs, fmt.Errorf("bar"))
	errs = multierror.Append(errs, fmt.Errorf("baz"))

	fmt.Printf("%v", errs)
	// Output:
	// 3 errors occurred:
	// foo
	// bar
	// baz
}

func ExampleKeyed() {
	var errs error

	errs = multierror.Append(errs, multierror.Keyed("k1", fmt.Errorf("foo")))
	errs = multierror.Append(errs, multierror.Keyed("k2", fmt.Errorf("foo")))
	errs = multierror.Append(errs, multierror.Keyed("k3", fmt.Errorf("bar")))

	fmt.Printf("%v", errs)
	// Output:
	// 3 errors occurred:
	// bar (k3)
	// foo (k1, k2)
}
