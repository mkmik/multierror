package multierror_test

import (
	"fmt"
	"strings"
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

func ExampleUnfold() {
	var errs error

	errs = multierror.Append(errs, fmt.Errorf("foo"))
	errs = multierror.Append(errs, fmt.Errorf("bar"))
	errs = multierror.Append(errs, fmt.Errorf("baz"))

	fmt.Printf("%q", multierror.Unfold(errs))
	// Output:
	// ["foo" "bar" "baz"]
}

func TestUnfoldSingleton(t *testing.T) {
	errs := multierror.Unfold(fmt.Errorf("foo"))
	if got, want := len(errs), 1; got != want {
		t.Fatalf("got: %d, want: %d", got, want)
	}
	if got, want := errs[0].Error(), "foo"; got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func ExampleUniq() {
	var errs []error

	errs = append(errs, fmt.Errorf("foo"))
	errs = append(errs, fmt.Errorf("bar"))
	errs = append(errs, fmt.Errorf("foo"))

	fmt.Printf("%q", multierror.Uniq(errs))
	// Output:
	// ["foo repeated 2 times" "bar"]
}

func ExampleTagged() {
	var errs error

	errs = multierror.Append(errs, multierror.Tagged("k1", fmt.Errorf("foo")))
	errs = multierror.Append(errs, multierror.Tagged("k2", fmt.Errorf("foo")))
	errs = multierror.Append(errs, multierror.Tagged("k3", fmt.Errorf("bar")))

	fmt.Printf("%v", errs)
	// Output:
	// 3 errors occurred:
	// foo (k1)
	// foo (k2)
	// bar (k3)
}

func ExampleTagged_uniq() {
	var errs []error

	errs = append(errs, multierror.Tagged("k1", fmt.Errorf("foo")))
	errs = append(errs, multierror.Tagged("k2", fmt.Errorf("foo")))
	errs = append(errs, multierror.Tagged("k3", fmt.Errorf("bar")))

	fmt.Printf("%q", multierror.Uniq(errs))
	// Output:
	// ["foo (k1, k2)" "bar (k3)"]
}

func ExampleFormatter() {
	var errs error

	errs = multierror.Append(errs, multierror.Tagged("k1", fmt.Errorf("foo")))
	errs = multierror.Append(errs, multierror.Tagged("k2", fmt.Errorf("foo")))
	errs = multierror.Append(errs, multierror.Tagged("k3", fmt.Errorf("bar")))

	errs = multierror.Append(nil, multierror.Uniq(multierror.Unfold(errs))...)

	errs = multierror.WithFormatter(errs, func(errs []string) string {
		return strings.Join(errs, "; ")
	})

	fmt.Printf("%v", errs)
	// Output:
	// foo (k1, k2); bar (k3)
}
