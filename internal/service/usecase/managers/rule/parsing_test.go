package manager

import (
	"context"
	"errors"
	"testing"
	"time"

	eschedule "jcalendar/internal/service/entity/schedule"

	"github.com/stretchr/testify/require"
)

func TestHandleRule(t *testing.T) {
	var (
		ctx = context.Background()

		eventFrom = time.Date(2022, 1, 3, 0, 0, 0, 0, time.Local)
	)

	r := NewRuleManager(ctx)
	table := []*struct {
		testSubTittle     string
		rule              string
		expectedSchedules []*eschedule.EventSchedule
		expectedError     error
	}{
		{
			testSubTittle:     "empty rule",
			rule:              "",
			expectedSchedules: nil,
			expectedError:     ErrInvalidRuleExpr,
		},
		{
			testSubTittle:     "invalid rule",
			rule:              "SCHEDULER_MODE=COMMON,ENDING_MODE=NONE,INTERVAL=1,IS_REGULAR=TRUE,SHIFT=DAILY",
			expectedSchedules: nil,
			expectedError:     ErrInvalidRuleExpr,
		},
		{
			testSubTittle:     "invalid rule parts",
			rule:              "SCHEDULER_MODE-COMMON;ENDING_MODE-NONE;INTERVAL-1;IS_REGULAR-TRUE;SHIFT-DAILY",
			expectedSchedules: nil,
			expectedError:     ErrInvalidPartLen,
		},
		{
			testSubTittle:     "invalid Interval part value",
			rule:              "SCHEDULER_MODE=COMMON;ENDING_MODE=NONE;INTERVAL=A;IS_REGULAR=TRUE;SHIFT=DAILY",
			expectedSchedules: nil,
			expectedError:     ErrInvalidIntervalPartVal,
		},
		{
			testSubTittle:     "invalid isRegular part value",
			rule:              "SCHEDULER_MODE=COMMON;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=AAAA;SHIFT=DAILY",
			expectedSchedules: nil,
			expectedError:     ErrInvalidIsRegularPartVal,
		},
		{
			testSubTittle:     "invalid endOccurrence part value",
			rule:              "SCHEDULER_MODE=CUSTOM;ENDING_MODE=DATA;END_OCCURRENCE=AAAA;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=WEEKLY;CUSTOM_DAY_LIST=1,2,3;",
			expectedSchedules: nil,
			expectedError:     ErrInvalidEndOccurrencePartVal,
		},
		{
			testSubTittle: "valid common rule",
			rule:          "SCHEDULER_MODE=COMMON;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=DAILY",
			expectedSchedules: []*eschedule.EventSchedule{
				{
					BeginOccurrence: eventFrom,
					SchedulerType:   "COMMON",
					EndingMode:      "NONE",
					IntervalVal:     1,
					IsRegular:       true,
					Shift:           "DAILY",
				},
			},
			expectedError: nil,
		},
		{
			testSubTittle: "valid custom rule",
			rule:          "SCHEDULER_MODE=CUSTOM;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=WEEKLY;CUSTOM_DAY_LIST=1,2,3",
			expectedSchedules: []*eschedule.EventSchedule{
				{
					BeginOccurrence: eventFrom,
					SchedulerType:   "CUSTOM",
					EndingMode:      "NONE",
					IntervalVal:     1,
					IsRegular:       true,
					Shift:           "WEEKLY",
				},
				{
					BeginOccurrence: eventFrom.Add(24 * time.Hour),
					SchedulerType:   "CUSTOM",
					EndingMode:      "NONE",
					IntervalVal:     1,
					IsRegular:       true,
					Shift:           "WEEKLY",
				},
				{
					BeginOccurrence: eventFrom.Add(48 * time.Hour),
					SchedulerType:   "CUSTOM",
					EndingMode:      "NONE",
					IntervalVal:     1,
					IsRegular:       true,
					Shift:           "WEEKLY",
				},
			},
			expectedError: nil,
		},
	}

	for _, row := range table {
		t.Run(row.testSubTittle, func(t *testing.T) {
			schs, err := r.HandleRule(ctx, eventFrom, row.rule)
			require.True(t, errors.Is(err, row.expectedError))
			require.Equal(t, row.expectedSchedules, schs)
		})
	}
}
