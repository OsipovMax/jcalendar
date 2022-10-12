package scheduler

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	cmdscheduler "jcalendar/internal/service/usecase/commands/scheduler"
)

const (
	customSchedulerMode = "CUSTOM"

	dailyShiftKey   = "DAILY"
	weeklyShiftKey  = "WEEKLY"
	monthlyShiftKey = "MONTHLY"
	yearlyShiftKey  = "YEARLY"

	schedulerModeKey = "SCHEDULER_MODE"
	endingModeKey    = "ENDING_MODE"
	intervalKey      = "INTERVAL"
	dayModeKey       = "IS_EACH_DAY"
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
SCHEDULER_MODE=CUSTOM;ENDING_MODE=NONE;INTERVAL=1;IS_EACH_DAY=TRUE;SHIFT=WEEKLY;CUSTOM_DAY_LIST=1,2,3;
ENDING_MODE=NONE or DATE or REPEAT_COUNT
INTERVAL=1+
*/

type EventScheduler struct {
	EventSchedulerCmdHandler *cmdscheduler.CreateEventScheduleCommandHandler
}

func NewEventScheduler(_ context.Context, handler *cmdscheduler.CreateEventScheduleCommandHandler) *EventScheduler {
	return &EventScheduler{
		EventSchedulerCmdHandler: handler,
	}
}

func (e *EventScheduler) HandleAndSaveSchedules(ctx context.Context, event *eevent.Event) error {
	_, err := e.handleRule(ctx, event)
	if err != nil {
		return err
	}

	//TODO: save schedules

	return nil
}

func (e *EventScheduler) handleRule(ctx context.Context, event *eevent.Event) ([]cmdscheduler.CreateEventScheduleCommand, error) {
	parts := strings.Split(event.ScheduleRule, ";")

	if len(parts) == 0 {
		return nil, errors.New("invalid rule expr")
	}

	return e.tokenize(ctx, event.ID, event.From, parts)
}

func (e *EventScheduler) tokenize(
	_ context.Context,
	eventID uint,
	eventStartedTimestamp time.Time,
	parts []string,
) ([]cmdscheduler.CreateEventScheduleCommand, error) {
	var (
		cmd         = cmdscheduler.CreateEventScheduleCommand{EventID: eventID, BeginOccurrence: eventStartedTimestamp}
		strDaysList string
	)

	for idx := range parts {
		part := strings.Split(parts[idx], "=")

		if len(part) < 2 {
			return nil, fmt.Errorf("invalid part - %v", part)
		}

		switch part[0] {
		case schedulerModeKey:
			cmd.SchedulerType = part[1]
		case endingModeKey:
			cmd.EndingMode = part[1]
		case intervalKey:
			intervalVal, err := strconv.Atoi(part[1])
			if err != nil {
				return nil, errors.New("invalid interval part")
			}

			cmd.IntervalVal = intervalVal
		case dayModeKey:
			isEachDay, err := strconv.ParseBool(part[1])
			if err != nil {
				return nil, errors.New("invalid isEachDay part")
			}

			cmd.IsEachDay = isEachDay
		case EndOccurrenceKey:
			t, err := time.Parse(time.RFC3339, part[1])
			if err != nil {
				return nil, errors.New("invalid EndOccurrence part")
			}

			cmd.EndOccurrence = t
		case shiftKey:
			switch part[1] {
			case dailyShiftKey:
				cmd.Daily = true
			case weeklyShiftKey:
				cmd.Weekly = true
			case monthlyShiftKey:
				cmd.Monthly = true
			case yearlyShiftKey:
				cmd.Yearly = true
			}
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

	if cmd.SchedulerType == customSchedulerMode && len(strDaysList) != 0 {
		curWeekday := cmd.BeginOccurrence.Weekday()
		var cmdList []cmdscheduler.CreateEventScheduleCommand
		for _, strDay := range strings.Split(strDaysList, ",") {
			day, err := strconv.Atoi(strDay)
			if err != nil {
				return nil, err
			}

			delta := (day - int(curWeekday) + 7) % 7

			cmd.BeginOccurrence.Add(time.Duration(delta * 24))
			cmdList = append(cmdList, cmd)
		}

		return cmdList, nil
	}

	return []cmdscheduler.CreateEventScheduleCommand{cmd}, nil
}
