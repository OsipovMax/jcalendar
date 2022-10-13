package manager

import (
	"context"
	"reflect"
	"testing"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	eschedule "jcalendar/internal/service/entity/schedule"

	"github.com/stretchr/testify/require"
)

func TestExtendWithScheduledEvents(t *testing.T) {
	var (
		ctx  = context.Background()
		from = time.Date(2022, 1, 1, 12, 30, 0, 0, time.Local)
		till = time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local)
	)

	e := NewEventManager(ctx, nil)
	table := []*struct {
		TestSubTittle  string
		From, Till     time.Time
		Events         []*eevent.Event
		ExpectedEvents []*eevent.Event
	}{
		{
			TestSubTittle: "unrepeated event in interval",
			From:          from,
			Till:          till,
			Events: []*eevent.Event{
				{
					From:     from.Add(3 * time.Minute),
					Till:     till.Add(-3 * time.Minute),
					Details:  "details",
					IsRepeat: false,
				},
			},
			ExpectedEvents: []*eevent.Event{
				{
					From:     from.Add(3 * time.Minute),
					Till:     till.Add(-3 * time.Minute),
					Details:  "details",
					IsRepeat: false,
				},
			},
		},
		{
			TestSubTittle: "unrepeated event not in interval",
			From:          from,
			Till:          till,
			Events: []*eevent.Event{
				{
					From:     from.Add(1 * time.Hour),
					Till:     till.Add(3 * time.Hour),
					Details:  "details",
					IsRepeat: false,
				},
			},
			ExpectedEvents: []*eevent.Event{},
		},
		{
			TestSubTittle: "repeated event + and scheduled event in interval (Inf mode)",
			From:          from,
			Till:          from.Add(45*time.Minute).AddDate(0, 0, 2),
			Events: []*eevent.Event{
				{
					From:    from,
					Till:    till,
					Details: "details",
					EventSchedules: []*eschedule.EventSchedule{
						{
							BeginOccurrence: from,
							IntervalVal:     2,
							EndingMode:      "NONE",
							Shift:           "DAILY",
							SchedulerType:   "CUSTOM",
						},
					},
					IsRepeat: true,
				},
			},
			ExpectedEvents: []*eevent.Event{
				{
					From:     from,
					Till:     till,
					Details:  "details",
					IsRepeat: true,
				},
				{
					From:     from.AddDate(0, 0, 2),
					Till:     till.AddDate(0, 0, 2),
					Details:  "details",
					IsRepeat: true,
				},
			},
		},
		{
			TestSubTittle: "scheduled event in interval (Inf mode)",
			From:          from,
			Till:          from.Add(45*time.Minute).AddDate(0, 0, 2),
			Events: []*eevent.Event{
				{
					From:    from.Add(-time.Hour),
					Till:    till,
					Details: "details",
					EventSchedules: []*eschedule.EventSchedule{
						{
							BeginOccurrence: from.Add(-time.Hour),
							IntervalVal:     2,
							EndingMode:      "NONE",
							Shift:           "DAILY",
							SchedulerType:   "CUSTOM",
						},
					},
					IsRepeat: true,
				},
			},
			ExpectedEvents: []*eevent.Event{
				{
					From:     from.Add(-time.Hour).AddDate(0, 0, 2),
					Till:     till.AddDate(0, 0, 2),
					Details:  "details",
					IsRepeat: true,
				},
			},
		},
		{
			TestSubTittle: "repeated event in interval (DATA MODE)",
			From:          from,
			Till:          till.AddDate(0, 0, 2),
			Events: []*eevent.Event{
				{
					From:    from,
					Till:    till,
					Details: "details",
					EventSchedules: []*eschedule.EventSchedule{
						{
							BeginOccurrence: from,
							EndOccurrence:   from.AddDate(0, 0, 1),
							IntervalVal:     2,
							EndingMode:      "NONE",
							Shift:           "DAILY",
							SchedulerType:   "CUSTOM",
						},
					},
					IsRepeat: true,
				},
			},
			ExpectedEvents: []*eevent.Event{
				{
					From:     from,
					Till:     till,
					Details:  "details",
					IsRepeat: true,
				},
			},
		},
	}

	for _, row := range table {
		t.Run(row.TestSubTittle, func(t *testing.T) {
			actualEvs, err := e.extendWithScheduledEvents(ctx, row.Events, row.From, row.Till)
			require.NoError(t, err)

			require.Equal(t, len(row.ExpectedEvents), len(actualEvs))
			for idx := range actualEvs {
				actualEvs[idx].EventSchedules = nil
			}

			require.True(t, reflect.DeepEqual(row.ExpectedEvents, actualEvs))
		})
	}
}
