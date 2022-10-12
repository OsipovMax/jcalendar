package event

import (
	"context"
	"errors"
	"time"
)

type GetEventsInIntervalQuery struct {
	UserID     uint
	From, Till time.Time
}

func NewGetEventsInIntervalQuery(
	_ context.Context,
	id uint,
	from, till time.Time,
) (*GetEventsInIntervalQuery, error) {
	if id <= 0 {
		return nil, errors.New("non-positive id value")
	}

	return &GetEventsInIntervalQuery{
		UserID: id,
		From:   from,
		Till:   till,
	}, nil
}
