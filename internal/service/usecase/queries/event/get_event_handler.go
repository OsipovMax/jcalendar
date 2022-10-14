package event

import (
	"context"

	eevent "jcalendar/internal/service/entity/event"
)

type Getter interface {
	GetEventByID(ctx context.Context, id uint) (*eevent.Event, error)
}

type GetEventQueryHandler struct {
	getter Getter
}

func NewGetEventQueryHandler(getter Getter) GetEventQueryHandler {
	return GetEventQueryHandler{getter: getter}
}

func (ch *GetEventQueryHandler) Handle(ctx context.Context, query *GetEventQuery) (*eevent.Event, error) {
	return ch.getter.GetEventByID(ctx, query.EventID)
}
