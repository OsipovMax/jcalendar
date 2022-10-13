package manager

import (
	"context"

	qrevent "jcalendar/internal/service/usecase/queries/event"
)

const (
	customSchedulerMode = "CUSTOM"

	schedulerModeKey = "SCHEDULER_MODE"
	endingModeKey    = "ENDING_MODE"
	intervalKey      = "INTERVAL"
	dayModeKey       = "IS_REGULAR"
	shiftKey         = "SHIFT"
	EndOccurrenceKey = "END_OCCURRENCE"
	CustomDayListKey = "CUSTOM_DAY_LIST"
)

const (
	dailyShiftKey   = "DAILY"
	weeklyShiftKey  = "WEEKLY"
	monthlyShiftKey = "MONTHLY"
	yearlyShiftKey  = "YEARLY"
)

type EventManager struct {
	EventsInIntervalQueryHandler *qrevent.GetEventsInIntervalQueryHandler
}

func NewEventManager(_ context.Context, handler *qrevent.GetEventsInIntervalQueryHandler) *EventManager {
	return &EventManager{
		EventsInIntervalQueryHandler: handler,
	}
}
