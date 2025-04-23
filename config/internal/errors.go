package internal

import (
	"fmt"

	"github.com/mandelsoft/goutils/errors"
)

const KIND_CONFIGTYPE = "config type"

////////////////////////////////////////////////////////////////////////////////

type noContextError struct {
	name string
}

func (e *noContextError) Error() string {
	return fmt.Sprintf("unknown context %q", e.name)
}

func ErrNoContext(name string) error {
	return &noContextError{name}
}

func IsErrNoContext(err error) bool {
	return errors.IsA(err, &noContextError{})
}

func IsErrConfigNotApplicable(err error) bool {
	return errors.IsErrUnknownKind(err, KIND_CONFIGTYPE) || IsErrNoContext(err)
}
