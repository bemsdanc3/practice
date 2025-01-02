package errs

import "errors"

var (
	ErrTicketAlreadyRefunded = errors.New("ticket already refunded")
	ErrTicketNotFound        = errors.New("ticket not found")
	ErrTicketAlreadySold     = errors.New("ticket already sold")
)
