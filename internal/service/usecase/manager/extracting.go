package manager

import (
	"context"
	"fmt"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	qrevent "jcalendar/internal/service/usecase/queries/event"
)

func (e *EventManager) GetEventsInInterval(ctx context.Context, userID uint, from, till string) ([]*eevent.Event, error) {
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
		if !ev.IsRepeat && ev.From.After(ft) { // || ending_mode = "DATA"
			fullEventsList = append(fullEventsList, ev)
			continue
		}

		for _, sch := range ev.EventSchedules {
			var (
				timestamp     = sch.BeginOccurrence
				eventDuration = ev.Till.Sub(ev.From)
			)

			if timestamp.Equal(ft) || timestamp.After(ft) && timestamp.Before(tt) {
				fullEventsList = append(fullEventsList, ev)
			}

			/*
				А не достигли ли мы конца интервала самого расписания ????
			*/

			for timestamp.Before(tt) {
				shift := sch.Shift
				switch shift {
				case dailyShiftKey:
					timestamp = timestamp.AddDate(0, 0, 1)
				case weeklyShiftKey:
					timestamp = timestamp.AddDate(0, 0, 7)
				case monthlyShiftKey:
					timestamp = timestamp.AddDate(0, 1, 0)
				case yearlyShiftKey:
					timestamp = timestamp.AddDate(1, 0, 0)
				}

				if timestamp.Before(tt) && timestamp.Add(eventDuration).Before(tt) {
					eventCopy := copyEvent(ev)
					ev.From = timestamp
					ev.Till = ev.From.Add(eventDuration)
					fullEventsList = append(fullEventsList, eventCopy)
				}
			}
		}
	}

	return fullEventsList, nil
}
