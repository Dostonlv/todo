package structs

import "errors"

var (
	ErrBadRequest      = errors.New("bad request")
	ErrNoRowsAffected  = errors.New("no rows affected")
	ErrNotFound        = errors.New("not found")
	ErrUserBlocked     = errors.New("user blocked")
	ErrUniqueViolation = errors.New("unique Violation error")
	ErrWhiteList       = errors.New("account in whitelist")
)
