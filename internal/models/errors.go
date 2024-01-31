package models

import (
	"errors"
)

var (
	ErrNoRecord              = errors.New("models: no matching record found")
	ErrInvalidCredentials    = errors.New("models: invalid credentials")
	ErrDuplicateEmail        = errors.New("models: duplicate email")
	ErrDuplicateNotification = errors.New("models: duplicate notification")
	ErrDuplicateClass        = errors.New("models: duplicate class")
)
