package manager

import (
	"errors"
)

var (
	ErrEmptyUserIDsList      = errors.New("empty IDs user list")
	ErrInvalidWindowDuration = errors.New("invalid window duration")
	ErrInvalidFromVal        = errors.New("invalid parsing from timestamp string")
	ErrInvalidTillVal        = errors.New("invalid converting till timestamp string")
	ErrInvalidCreatingQuery  = errors.New("invalid creating query")
	ErrInvalidExecutingQuery = errors.New("invalid executing query")
)
