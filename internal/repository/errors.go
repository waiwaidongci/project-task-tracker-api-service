package repository

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func notFound(format string, args ...any) error {
	return fmt.Errorf(format+": %w", append(args, ErrNotFound)...)
}
