package manager

import "errors"

var (
	ErrInvalidRuleExpr             = errors.New("invalid rule expr")
	ErrUnknownRuleToken            = errors.New("unknown rule token")
	ErrInvalidPartLen              = errors.New("invalid part len")
	ErrInvalidIntervalPartVal      = errors.New("invalid interval part value")
	ErrInvalidIsRegularPartVal     = errors.New("invalid isRegular part value")
	ErrInvalidEndOccurrencePartVal = errors.New("invalid EndOccurrence part value")
)
