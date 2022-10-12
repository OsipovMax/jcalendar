package schedule

import (
	"context"
	"time"
)

type EventSchedule struct {
	ID              uint `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BeginOccurrence time.Time
	EndOccurrence   time.Time
	EndingMode      string
	IntervalVal     int
	Shift           string
	IsRegular       bool
	SchedulerType   string
	EventID         uint
}

func NewEventsSchedule(
	_ context.Context,
	beginOccurrence, endOccurrence time.Time,
	endingMode, shift, schedulerType string,
	intervalVal int,
	isRegular bool,
	eventID uint,
) *EventSchedule {
	return &EventSchedule{
		BeginOccurrence: beginOccurrence,
		EndOccurrence:   endOccurrence,
		EndingMode:      endingMode,
		IntervalVal:     intervalVal,
		Shift:           shift,
		IsRegular:       isRegular,
		SchedulerType:   schedulerType,
		EventID:         eventID,
	}
}
