package manager

import (
	"container/heap"
	"context"
	"time"

	"github.com/sirupsen/logrus"

	eevent "jcalendar/internal/service/entity/event"
)

func (e *EventManager) GetClosestFreeWindow(ctx context.Context, userIDs []int, winSize string) (time.Time, time.Time, error) {
	if len(userIDs) == 0 {
		return time.Time{}, time.Time{}, ErrEmptyUserIDsList
	}

	now := time.Now()
	winDuration, err := time.ParseDuration(winSize)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t parse winSize param: %v", err)
		return time.Time{}, time.Time{}, ErrInvalidWindowDuration
	}

	intervals := make([]*Interval, 0)
	for _, userID := range userIDs {
		var evs []*eevent.Event
		evs, err = e.GetEventsInInterval(ctx, uint(userID), now.Format(time.RFC3339), now.AddDate(1, 0, 0).Format(time.RFC3339))
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		for _, ev := range evs {
			intervals = append(intervals, &Interval{From: ev.From, Till: ev.Till})
		}
	}

	from, till := e.getClosestFreeWindow(ctx, intervals, winDuration)

	return from, till, nil
}

func (e *EventManager) getClosestFreeWindow(_ context.Context, intervals []*Interval, winSize time.Duration) (time.Time, time.Time) {
	var intervalHeap IntervalHeap
	heap.Init(&intervalHeap)

	if len(intervals) == 0 {
		now := time.Now()
		return now, now.Add(winSize)
	}

	for _, interval := range intervals {
		heap.Push(&intervalHeap, interval)
	}

	curInterval := heap.Pop(&intervalHeap).(*Interval)

	now := time.Now()
	if curInterval.From.Sub(now) > winSize {
		return now, now.Add(winSize)
	}

	for intervalHeap.Len() > 0 {
		nextInterval := heap.Pop(&intervalHeap).(*Interval)

		if nextInterval.From.Sub(curInterval.Till) > winSize {
			break
		} else if nextInterval.Till.After(curInterval.Till) {
			curInterval = nextInterval
		}
	}

	return curInterval.Till, curInterval.Till.Add(winSize)
}
