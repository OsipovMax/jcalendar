package event

import (
	"context"
	"time"

	eevent "jcalendar/internal/service/entity/event"
)

type IntervalGetter interface {
	GetEventsInInterval(ctx context.Context, userID uint, from, till time.Time) ([]*eevent.Event, error)
}

type GetEventsInIntervalQueryHandler struct {
	getter IntervalGetter
}

func NewGetEventsInIntervalQueryHandler(getter IntervalGetter) *GetEventsInIntervalQueryHandler {
	return &GetEventsInIntervalQueryHandler{getter: getter}
}

func (ch *GetEventsInIntervalQueryHandler) Handle(ctx context.Context, query *GetEventsInIntervalQuery) ([]*eevent.Event, error) {
	es, err := ch.getter.GetEventsInInterval(ctx, query.UserID, query.From, query.Till)
	if err != nil {
		return nil, err
	}

	return es, nil
}
