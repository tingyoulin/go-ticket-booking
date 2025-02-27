package domain

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
	// ErrUnauthorized will throw if the user is not authorized
	ErrUnauthorized = errors.New("unauthorized (invalid token or password)")
	// ErrForbidden will throw if the user is forbidden to access the resource
	ErrForbidden = errors.New("forbidden (do not have permission)")
	// ErrFlightCanceled will throw if the flight is canceled
	ErrFlightCanceled = errors.New("flight is canceled")
	// ErrFlightNoAvailableSeats will throw if the flight has no available seats
	ErrFlightNoAvailableSeats = errors.New("flight has no available seats")
)
