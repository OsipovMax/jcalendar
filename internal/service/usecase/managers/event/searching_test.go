package manager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetClosestFreeWindow(t *testing.T) {
	var (
		ctx         = context.Background()
		winDuration = 15 * time.Minute
		from        = time.Date(2022, 1, 1, 12, 30, 0, 0, time.Local)
		till        = time.Date(2022, 1, 1, 13, 0, 0, 0, time.Local)
	)

	e := NewEventManager(ctx, nil)
	table := []*struct {
		testSubTittle              string
		intervals                  []*Interval
		winDuration                time.Duration
		expectedFrom, expectedTill time.Time
	}{
		{
			testSubTittle: "single time interval",
			intervals:     []*Interval{{From: from, Till: till}},
			winDuration:   winDuration,
			expectedFrom:  till,
			expectedTill:  till.Add(winDuration),
		},
		{
			testSubTittle: "F1___T1_F2___T2",
			intervals:     []*Interval{{From: from, Till: till}, {From: till.Add(winDuration / 2), Till: till.Add(winDuration * 3)}},
			winDuration:   winDuration,
			expectedFrom:  till.Add(winDuration * 3),
			expectedTill:  till.Add(winDuration * 4),
		},
		{
			testSubTittle: "F1___T1___F2___T2",
			intervals:     []*Interval{{From: from, Till: till}, {From: till.Add(winDuration * 2), Till: till.Add(winDuration * 3)}},
			winDuration:   winDuration,
			expectedFrom:  till,
			expectedTill:  till.Add(winDuration),
		},
		{
			testSubTittle: "F1___F2___T2___T1",
			intervals:     []*Interval{{From: from, Till: till}, {From: from.Add(winDuration), Till: from.Add(winDuration + 5)}},
			winDuration:   winDuration,
			expectedFrom:  till,
			expectedTill:  till.Add(winDuration),
		},
		{
			testSubTittle: "F1___F2___T1___T2",
			intervals:     []*Interval{{From: from, Till: till}, {From: from.Add(winDuration), Till: till.Add(winDuration)}},
			winDuration:   winDuration,
			expectedFrom:  till.Add(winDuration),
			expectedTill:  till.Add(winDuration * 2),
		},
		{
			testSubTittle: "F1___T1___F2___T3___T2___F3",
			intervals: []*Interval{
				{From: from, Till: till},
				{From: till.Add(winDuration * 2), Till: till.Add(winDuration * 5)},
				{From: till.Add(winDuration * 3), Till: till.Add(winDuration * 7)},
			},
			winDuration:  winDuration,
			expectedFrom: till,
			expectedTill: till.Add(winDuration),
		},
	}

	for _, row := range table {
		t.Run(row.testSubTittle, func(t *testing.T) {
			actualFrom, actualTill := e.getClosestFreeWindow(ctx, row.intervals, row.winDuration)
			require.Equal(t, row.expectedFrom, actualFrom)
			require.Equal(t, row.expectedTill, actualTill)
		})
	}
}
