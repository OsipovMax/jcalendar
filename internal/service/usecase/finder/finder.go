package finder

import (
	"container/heap"
	"context"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	"jcalendar/internal/service/usecase/extractor"
)

type Interval struct {
	From, Till time.Time
}

type Finder struct {
	Extractor *extractor.EventExtractor
}

func NewFinder(_ context.Context, extractor *extractor.EventExtractor) *Finder {
	return &Finder{Extractor: extractor}
}

func (f *Finder) GetClosestFreeWindow(ctx context.Context, userIDs []int, winSize string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	winDuration, err := time.ParseDuration(winSize)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	intervals := make([]*Interval, 0)
	for _, userID := range userIDs {
		var evs []*eevent.Event
		evs, err = f.Extractor.GetEventsInInterval(ctx, uint(userID), now.Format(time.RFC3339), now.AddDate(1, 0, 0).Format(time.RFC3339))
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		for _, ev := range evs {
			intervals = append(intervals, &Interval{From: ev.From, Till: ev.Till})
		}
	}

	from, till := f.getClosestFreeWindow(ctx, intervals, winDuration)

	return from, till, nil
}

func (f *Finder) getClosestFreeWindow(_ context.Context, intervals []*Interval, winSize time.Duration) (time.Time, time.Time) {
	var intervalHeap IntervalHeap
	heap.Init(&intervalHeap)

	if len(intervals) == 0 {
		now := time.Now().UTC()
		return now, now.Add(winSize)
	}

	for _, interval := range intervals {
		heap.Push(&intervalHeap, interval)
	}

	prevInterval := heap.Pop(&intervalHeap).(*Interval)
	for intervalHeap.Len() > 0 {
		interval := heap.Pop(&intervalHeap).(*Interval)

		if prevInterval.Till.Before(interval.From) && interval.From.Sub(prevInterval.Till) > winSize {
			break
		} else {
			if interval.Till.After(prevInterval.Till) {
				prevInterval = interval
			}
		}
	}

	return prevInterval.Till, prevInterval.Till.Add(winSize)
}
