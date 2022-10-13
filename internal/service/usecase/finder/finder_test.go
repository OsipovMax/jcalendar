package finder

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

	f := NewFinder(ctx, nil)
	table := []*struct {
		TestSubTittle              string
		Intervals                  []*Interval
		WinDuration                time.Duration
		ExpectedFrom, ExpectedTill time.Time
	}{
		//{
		//	TestSubTittle: "empty interval list",
		//	Intervals:     nil,
		//	WinDuration:   winDuration,
		//},
		{
			TestSubTittle: "single time interval",
			Intervals:     []*Interval{{From: from, Till: till}},
			WinDuration:   winDuration,
			ExpectedFrom:  till,
			ExpectedTill:  till.Add(winDuration),
		},
		{
			TestSubTittle: "single time interval",
			Intervals:     []*Interval{{From: from, Till: till}},
			WinDuration:   winDuration,
			ExpectedFrom:  till,
			ExpectedTill:  till.Add(winDuration),
		},
		{
			TestSubTittle: "F1___T1___F2___T2",
			Intervals:     []*Interval{{From: from, Till: till}, {From: till.Add(winDuration * 2), Till: till.Add(winDuration * 3)}},
			WinDuration:   winDuration,
			ExpectedFrom:  till,
			ExpectedTill:  till.Add(winDuration),
		},
		{
			TestSubTittle: "F1___F2___T2___T1",
			Intervals:     []*Interval{{From: from, Till: till}, {From: from.Add(winDuration), Till: from.Add(winDuration + 5)}},
			WinDuration:   winDuration,
			ExpectedFrom:  till,
			ExpectedTill:  till.Add(winDuration),
		},
		{
			TestSubTittle: "F1___F2___T1___T2",
			Intervals:     []*Interval{{From: from, Till: till}, {From: from.Add(winDuration), Till: till.Add(winDuration)}},
			WinDuration:   winDuration,
			ExpectedFrom:  till.Add(winDuration),
			ExpectedTill:  till.Add(winDuration * 2),
		},
		{
			TestSubTittle: "F1___T1___F2___T3___T2___F3",
			Intervals: []*Interval{
				{From: from, Till: till},
				{From: till.Add(winDuration * 2), Till: till.Add(winDuration * 5)},
				{From: till.Add(winDuration * 3), Till: till.Add(winDuration * 7)},
			},
			WinDuration:  winDuration,
			ExpectedFrom: till,
			ExpectedTill: till.Add(winDuration),
		},
	}

	for _, row := range table {
		t.Run(row.TestSubTittle, func(t *testing.T) {
			actualFrom, actualTill := f.getClosestFreeWindow(ctx, row.Intervals, row.WinDuration)
			require.Equal(t, row.ExpectedFrom, actualFrom)
			require.Equal(t, row.ExpectedTill, actualTill)
		})
	}
}
