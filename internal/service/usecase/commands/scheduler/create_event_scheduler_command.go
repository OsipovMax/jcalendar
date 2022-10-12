package scheduler

import (
	"context"
	"errors"
	"time"
)

type CreateEventScheduleCommand struct {
	BeginOccurrence time.Time
	EndOccurrence   time.Time
	EndingMode      string
	IntervalVal     int
	Shift           string
	IsRegular       bool
	SchedulerType   string
	EventID         uint
}

func NewCreateEventScheduleCommand(
	_ context.Context,
	beginOccurrence, endOccurrence time.Time,
	endingMode, schedulerType string,
	intervalVal int,
	shift string,
	isRegular bool,
	eventID uint,
) (*CreateEventScheduleCommand, error) {
	if endingMode == "" {
		return nil, errors.New("missing endingMode value")
	}

	if schedulerType == "" {
		return nil, errors.New("missing schedulerType value")
	}

	return &CreateEventScheduleCommand{
		BeginOccurrence: beginOccurrence,
		EndOccurrence:   endOccurrence,
		EndingMode:      endingMode,
		IntervalVal:     intervalVal,
		Shift:           shift,
		IsRegular:       isRegular,
		SchedulerType:   schedulerType,
		EventID:         eventID,
	}, nil
}
