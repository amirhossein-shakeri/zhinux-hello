package greeting

import "errors"

var (
	ErrInvalidPolicy         = errors.New("invalid greeting policy")
	ErrNameEmpty             = errors.New("name cannot be empty")
	ErrNameTooLong           = errors.New("name exceeds max length")
	ErrMessageEmpty          = errors.New("message cannot be empty")
	ErrMessageTooLong        = errors.New("message exceeds max length")
	ErrInvalidUTF8           = errors.New("value must be valid UTF-8")
	ErrStreamCountOutOfRange = errors.New("stream count is out of allowed range")
)
