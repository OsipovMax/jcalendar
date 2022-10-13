package manager

import (
	"context"
	"fmt"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	qrevent "jcalendar/internal/service/usecase/queries/event"
)

const (
	dataEndingMode = "DATA"
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

	return e.extendWithScheduledEvents(ctx, evs, ft, tt)
}

func (e *EventManager) extendWithScheduledEvents(_ context.Context, evs []*eevent.Event, ft, tt time.Time) ([]*eevent.Event, error) {
	fullEventsList := make([]*eevent.Event, 0)
	for _, ev := range evs {
		if !ev.IsRepeat && (ev.From.Equal(ft) || ev.From.After(ft)) && ev.Till.Before(tt) {
			fullEventsList = append(fullEventsList, ev)
			continue
		}

		for _, sch := range ev.EventSchedules {
			var (
				timestamp     = sch.BeginOccurrence
				eventDuration = ev.Till.Sub(ev.From)
			)

			if (timestamp.Equal(ft) || timestamp.After(ft)) && timestamp.Before(tt) {
				fmt.Println("!!!!!!!")
				fullEventsList = append(fullEventsList, ev)
			}

			endTimeStamp := tt
			if sch.EndingMode == dataEndingMode {
				endTimeStamp = sch.EndOccurrence
			}

			for timestamp.Before(endTimeStamp) {
				shift := sch.Shift
				switch shift {
				case dailyShiftKey:
					timestamp = timestamp.AddDate(0, 0, 1*sch.IntervalVal)
				case weeklyShiftKey:
					timestamp = timestamp.AddDate(0, 0, 7*sch.IntervalVal)
				case monthlyShiftKey:
					timestamp = timestamp.AddDate(0, 1*sch.IntervalVal, 0)
				case yearlyShiftKey:
					timestamp = timestamp.AddDate(1*sch.IntervalVal, 0, 0)
				}
				fmt.Println(timestamp)
				if timestamp.Before(tt) && timestamp.Add(eventDuration).Before(tt) {
					eventCopy := copyEvent(ev)
					eventCopy.From = timestamp
					eventCopy.Till = eventCopy.From.Add(eventDuration)
					fullEventsList = append(fullEventsList, eventCopy)
				}
			}
		}
	}

	return fullEventsList, nil
}
