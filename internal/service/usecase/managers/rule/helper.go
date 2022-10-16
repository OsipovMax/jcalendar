package manager

import eschedule "jcalendar/internal/service/entity/schedule"

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
