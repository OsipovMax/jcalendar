package manager

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	eschedule "jcalendar/internal/service/entity/schedule"
)

func (e *EventManager) HandleRule(
	ctx context.Context,
	eventFrom time.Time,
	eventScheduleRule string,
) ([]*eschedule.EventSchedule, error) {
	parts := strings.Split(eventScheduleRule, ";")

	if len(parts) == 0 {
		return nil, errors.New("invalid rule expr")
	}

	return e.tokenize(ctx, eventFrom, parts)
}

func (e *EventManager) tokenize(
	_ context.Context,
	eventStartedTimestamp time.Time,
	parts []string,
) ([]*eschedule.EventSchedule, error) {
	var (
		schedule    = &eschedule.EventSchedule{BeginOccurrence: eventStartedTimestamp}
		strDaysList string
	)

	for idx := range parts {
		part := strings.Split(parts[idx], "=")

		if len(part) < 2 {
			return nil, fmt.Errorf("invalid part - %v", part)
		}

		switch part[0] {
		case schedulerModeKey:
			schedule.SchedulerType = part[1]
		case endingModeKey:
			schedule.EndingMode = part[1]
		case intervalKey:
			intervalVal, err := strconv.Atoi(part[1])
			if err != nil {
				return nil, errors.New("invalid interval part")
			}

			schedule.IntervalVal = intervalVal
		case dayModeKey:
			isRegular, err := strconv.ParseBool(part[1])
			if err != nil {
				return nil, errors.New("invalid isEachDay part")
			}

			schedule.IsRegular = isRegular
		case EndOccurrenceKey:
			t, err := time.Parse(time.RFC3339, part[1])
			if err != nil {
				return nil, errors.New("invalid EndOccurrence part")
			}

			schedule.EndOccurrence = t
		case shiftKey:
			schedule.Shift = part[1]
		case CustomDayListKey:
			strDaysList = part[1]
		default:
			return nil, errors.New("invalid rule format")
		}
	}

	if schedule.SchedulerType == customSchedulerMode && len(strDaysList) != 0 {
		curWeekday := schedule.BeginOccurrence.Weekday()
		var scheduleList []*eschedule.EventSchedule
		for _, strDay := range strings.Split(strDaysList, ",") {
			day, err := strconv.Atoi(strDay)
			if err != nil {
				return nil, err
			}

			delta := (day - int(curWeekday) + 7*schedule.IntervalVal) % 7

			cSchedule := copySchedule(schedule)
			cSchedule.BeginOccurrence = cSchedule.BeginOccurrence.AddDate(0, 0, delta)
			scheduleList = append(scheduleList, cSchedule)
		}

		return scheduleList, nil
	}

	return []*eschedule.EventSchedule{schedule}, nil
}
