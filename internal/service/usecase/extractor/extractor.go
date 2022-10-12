package extractor

import (
	"context"
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

type Extractor struct {
	EventsInIntervalQueryHandler *qrevent.GetEventsInIntervalQueryHandler
}

func NewExtractor(_ context.Context, handler *qrevent.GetEventsInIntervalQueryHandler) *Extractor {
	return &Extractor{
		EventsInIntervalQueryHandler: handler,
	}
}

func (e *Extractor) GetEventsInInterval(ctx context.Context, userID uint, from, till time.Time) ([]*eevent.Event, error) {
	q, err := qrevent.NewGetEventsInIntervalQuery(ctx, userID, from, till)
	if err != nil {
		return nil, err
	}

	evs, err := e.EventsInIntervalQueryHandler.Handle(ctx, q)
	if err != nil {
		return nil, err
	}

	fullEventsList := make([]*eevent.Event, 0)
	for _, ev := range evs {
		if !ev.IsRepeat && ev.Till.Before(till) {
			fullEventsList = append(fullEventsList, ev)
			continue
		}

		if ev.From.After(from) && ev.From.Before(till) {
			fullEventsList = append(fullEventsList, ev)
		}

		var (
			curTimestamp  = ev.Schedule.BeginOccurrence
			eventDuration = ev.Till.Sub(ev.From)
		)

		/*
			А не достигли ли мы конца интервала самого расписания
		*/
		for curTimestamp.Before(till) {
			shift := ev.Schedule.SchedulerType // TODO SHIFT
			switch shift {
			case dailyShiftKey:
				curTimestamp.AddDate(0, 0, 1)
			case weeklyShiftKey:
				curTimestamp.AddDate(0, 0, 7)
			case monthlyShiftKey:
				curTimestamp.AddDate(0, 1, 0)
			case yearlyShiftKey:
				curTimestamp.AddDate(1, 0, 0)
			}

			if curTimestamp.Before(till) && curTimestamp.Add(eventDuration).Before(till) {
				eventCopy := copyEvent(ev)
				ev.From = curTimestamp
				ev.Till = ev.From.Add(eventDuration)
				fullEventsList = append(fullEventsList, eventCopy)
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
		Participants:    e.Participants,
		Invites:         e.Invites,
		Details:         e.Details,
		ScheduleRule:    e.ScheduleRule,
		Schedule:        e.Schedule,
		IsRepeat:        e.IsRepeat,
		IsPrivate:       e.IsPrivate,
	}
}
