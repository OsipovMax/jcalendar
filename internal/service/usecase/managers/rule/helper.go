package manager

import eschedule "jcalendar/internal/service/entity/schedule"

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
