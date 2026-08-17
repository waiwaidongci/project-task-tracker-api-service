package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func notFound(format string, args ...any) error {
	return fmt.Errorf(format+": %v", append(args, sql.ErrNoRows)...)
}
