package parser

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	eschedule "jcalendar/internal/service/entity/schedule"
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

/*
RULE_SYNTAX:
1. SCHEDULER_MODE=COMMON or CUSTOM
2. ENDING_MODE=NONE or DATE or REPEAT_COUNT
3. INTERVAL=1+
-------
SCHEDULER_MODE=COMMON;ENDING_MODE=NONE;INTERVAL=1;IS_EACH_DAY=TRUE;SHIFT=DAILY;

SCHEDULER_MODE=COMMON
ENDING_MODE=NONE
INTERVAL=1
-> Daily or Weekly or Monthly or Yearly
-------
SCHEDULER_MODE=CUSTOM;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=WEEKLY;CUSTOM_DAY_LIST=1,2,3;
ENDING_MODE=NONE or DATE or REPEAT_COUNT
INTERVAL=1+
*/

type ScheduleParser struct{}

func NewScheduleParser(_ context.Context) *ScheduleParser {
	return &ScheduleParser{}
}

//func (e *EventScheduler) HandleAndSaveSchedules(ctx context.Context, event *eevent.Event) error {
//	cmds, err := e.handleRule(ctx, event)
//	if err != nil {
//		return err
//	}
//
//	//TODO: replace on batch insert
//	for idx := range cmds {
//		_, err := e.EventSchedulerCmdHandler.Handle(ctx, &cmds[idx])
//	}
//
//	return nil
//}

func (e *ScheduleParser) HandleRule(
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

func (e *ScheduleParser) tokenize(
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

	/*

			3    5(now)  6
			3 - 5 = (-2 + 7) % 7 = 5
		    5 - 5 = (0 + 7) % 7 = 0
			6 - 5 = (1 + 7) % 7 = 1

	*/

	if schedule.SchedulerType == customSchedulerMode && len(strDaysList) != 0 {
		curWeekday := schedule.BeginOccurrence.Weekday()
		var scheduleList []*eschedule.EventSchedule
		for _, strDay := range strings.Split(strDaysList, ",") {
			day, err := strconv.Atoi(strDay)
			if err != nil {
				return nil, err
			}

			delta := (day - int(curWeekday) + 7) % 7

			cSchedule := copySchedule(schedule)
			cSchedule.BeginOccurrence = cSchedule.BeginOccurrence.AddDate(0, 0, delta)
			scheduleList = append(scheduleList, cSchedule)
		}

		return scheduleList, nil
	}

	return []*eschedule.EventSchedule{schedule}, nil
}

func copySchedule(s *eschedule.EventSchedule) *eschedule.EventSchedule {
	return &eschedule.EventSchedule{
		ID:              s.ID,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
		BeginOccurrence: s.BeginOccurrence,
		EndOccurrence:   s.EndOccurrence,
		EndingMode:      s.EndingMode,
		IntervalVal:     s.IntervalVal,
		Shift:           s.Shift,
		IsRegular:       s.IsRegular,
		SchedulerType:   s.SchedulerType,
		EventID:         s.EventID,
	}
}
