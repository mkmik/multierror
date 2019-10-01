// +build go1.13

package multierror

import (
	"errors"
	"fmt"
)

// unwrap wraps go 1.13 Unwrap method
func unwrap(err error) error {
	return errors.Unwrap(err)
}

func errorSuffix(format string, err error, a ...interface{}) error {
	return fmt.Errorf(fmt.Sprintf("%w %s", err, format), a...)
}
