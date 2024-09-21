package customerrors

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
)
