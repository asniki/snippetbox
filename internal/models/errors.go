package models

import "errors"

var (
	// ErrNoRows is returned when there is no row matches the query
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials is returned when user tries to login with an incorrect email address or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail is returned when user tries to signup with an email address that's already in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
