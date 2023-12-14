//go:build go1.20
// +build go1.20

package multierror_test

import (
	"errors"
	"fmt"

	"github.com/mkmik/multierror"
)

// multierror.Split works also with multierrors produced by the new stdlib errors.Join
func ExampleErrorsSplit() {
	err := errors.Join(fmt.Errorf("foo"), fmt.Errorf("bar"), fmt.Errorf("baz"))

	fmt.Printf("%q", multierror.Split(err))
	// Output:
	// ["foo" "bar" "baz"]
}

// We can combine the errors.Join way of joining errors with the extra features of multierror
func ExampleErrorsTransformer() {
	err := errors.Join(
		multierror.Tag("k1", fmt.Errorf("foo")),
		multierror.Tag("k2", fmt.Errorf("foo")),
		multierror.Tag("k3", fmt.Errorf("bar")),
	)

	err = multierror.Transform(err, multierror.Uniq)
	err = multierror.Format(err, multierror.InlineFormatter)

	fmt.Printf("%v", err)

	// Output:
	// foo (k1, k2); bar (k3)
}
