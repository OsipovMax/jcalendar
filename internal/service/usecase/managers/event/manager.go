package manager

import (
	"context"

	qrevent "jcalendar/internal/service/usecase/queries/event"
)

type EventManager struct {
	EventsInIntervalQueryHandler *qrevent.GetEventsInIntervalQueryHandler
}

func NewEventManager(_ context.Context, handler *qrevent.GetEventsInIntervalQueryHandler) *EventManager {
	return &EventManager{
		EventsInIntervalQueryHandler: handler,
	}
}
