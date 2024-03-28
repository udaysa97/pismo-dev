package errors

import "errors"

var (
	ErrNotExists           error = errors.New("not found in repository")
	ErrDuplicateNotAllowed error = errors.New("duplicate not allowed")
	ErrForeignKeyViolation error = errors.New("Foreign key violation")
	ErrSQLStatement        error = errors.New("query: %s, query data: %s")
	ErrSQLQuery            error = errors.New("SQL Statement error, statement: %s, error: %s, error_metadata: %s")
	ErrToSQLStatement      error = errors.New("unable to get SQL statement")
	ErrInvalidData         error = errors.New("Invalid %s: %v")
)
