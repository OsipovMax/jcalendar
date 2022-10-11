package event

import (
	"errors"
)

type Query struct {
	EventID uint
}

func NewQuery(id uint) (*Query, error) {
	if id <= 0 {
		return nil, errors.New("non-positive id value")
	}

	return &Query{
		EventID: id,
	}, nil
}
