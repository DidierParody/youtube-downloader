package domain

import "errors"

// Custom domain errors used across the application.
// These errors are returned by use cases and handled by HTTP handlers
// to produce consistent JSON error responses.
var (
	// ErrNotFound is returned when a requested resource does not exist.
	ErrNotFound = errors.New("resource not found")

	// ErrConflict is returned when a resource already exists (e.g., duplicate email/username).
	ErrConflict = errors.New("resource conflict")

	// ErrUnauthorized is returned when authentication fails or token is invalid.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when the user lacks permission.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidInput is returned when request validation fails.
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal is returned for unclassified internal failures.
	ErrInternal = errors.New("internal server error")
)
