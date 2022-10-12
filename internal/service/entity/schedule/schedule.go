package schedule

import "time"

type EventsSchedule struct {
	ID              uint `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BeginOccurrence time.Time
	EndOccurrence   time.Time
	EndingMode      string
	IntervalVal     int
	Daily           bool
	Weekly          bool
	Monthly         bool
	Yearly          bool
	SchedulerType   string
	EventID         uint
}
