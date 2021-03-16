package models

import "errors"

var (
	SameForumeError = errors.New("duplicate key value violates unique constraint \"forum_pkey\" (SQLSTATE 23505)")
	NoUser = errors.New("")
)
