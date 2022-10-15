package schedule

import (
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
