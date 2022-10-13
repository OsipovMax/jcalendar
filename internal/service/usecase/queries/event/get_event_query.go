package event

import (
	"context"
	"errors"
)

type GetEventQuery struct {
	EventID uint
}

func NewGetEventQuery(_ context.Context, id uint) (*GetEventQuery, error) {
	if id <= 0 {
		return nil, errors.New("non-positive id value")
	}

	return &GetEventQuery{
		EventID: id,
	}, nil
}
