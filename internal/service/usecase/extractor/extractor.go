package extractor

import (
	"context"
	"fmt"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	qrevent "jcalendar/internal/service/usecase/queries/event"
)

const (
	dailyShiftKey   = "DAILY"
	weeklyShiftKey  = "WEEKLY"
	monthlyShiftKey = "MONTHLY"
	yearlyShiftKey  = "YEARLY"
)

type EventExtractor struct {
	EventsInIntervalQueryHandler *qrevent.GetEventsInIntervalQueryHandler
}

func NewEventExtractor(_ context.Context, handler qrevent.GetEventsInIntervalQueryHandler) *EventExtractor {
	return &EventExtractor{
		EventsInIntervalQueryHandler: &handler,
	}
}

func (e *EventExtractor) GetEventsInInterval(ctx context.Context, userID uint, from, till string) ([]*eevent.Event, error) {
	ft, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid converting from data for getting events in interval: %w", err)
	}

	tt, err := time.Parse(time.RFC3339, till)
	if err != nil {
		return nil, fmt.Errorf("invalid converting till data for getting events in interval: %w", err)
	}

	q, err := qrevent.NewGetEventsInIntervalQuery(ctx, userID, ft, tt)
	if err != nil {
		return nil, err
	}

	evs, err := e.EventsInIntervalQueryHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	fullEventsList := make([]*eevent.Event, 0)
	for _, ev := range evs {
		if !ev.IsRepeat && ev.Till.Before(tt) {
			fullEventsList = append(fullEventsList, ev)
			continue
		}

		for _, sch := range ev.Schedule {
			if ev.From.After(ft) && ev.From.Before(tt) {
				fullEventsList = append(fullEventsList, ev)
			}

			var (
				curTimestamp  = sch.BeginOccurrence
				eventDuration = ev.Till.Sub(ev.From)
			)

			/*
				А не достигли ли мы конца интервала самого расписания ????
			*/
			for curTimestamp.Before(tt) {
				shift := sch.Shift
				switch shift {
				case dailyShiftKey:
					curTimestamp = curTimestamp.AddDate(0, 0, 1)
				case weeklyShiftKey:
					curTimestamp = curTimestamp.AddDate(0, 0, 7)
				case monthlyShiftKey:
					curTimestamp = curTimestamp.AddDate(0, 1, 0)
				case yearlyShiftKey:
					curTimestamp = curTimestamp.AddDate(1, 0, 0)
				}

				if curTimestamp.Before(tt) && curTimestamp.Add(eventDuration).Before(tt) {
					eventCopy := copyEvent(ev)
					ev.From = curTimestamp
					ev.Till = ev.From.Add(eventDuration)
					fullEventsList = append(fullEventsList, eventCopy)
				}
			}
		}
	}

	return fullEventsList, nil
}

func copyEvent(e *eevent.Event) *eevent.Event {
	return &eevent.Event{
		ID:              e.ID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		From:            e.From,
		Till:            e.Till,
		CreatorID:       e.CreatorID,
		Creator:         e.Creator,
		ParticipantsIDs: e.ParticipantsIDs,
		Users:           e.Users,
		Invites:         e.Invites,
		Details:         e.Details,
		ScheduleRule:    e.ScheduleRule,
		Schedule:        e.Schedule,
		IsRepeat:        e.IsRepeat,
		IsPrivate:       e.IsPrivate,
	}
}
