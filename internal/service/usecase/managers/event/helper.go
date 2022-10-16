package manager

import (
	"time"

	eevent "jcalendar/internal/service/entity/event"
)

type Interval struct {
	From, Till time.Time
}

type IntervalHeap []*Interval

func (h IntervalHeap) Len() int {
	return len(h)
}
func (h IntervalHeap) Less(i, j int) bool {
	return h[i].From.Before(h[j].From)
}
func (h IntervalHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntervalHeap) Push(x interface{}) {
	*h = append(*h, x.(*Interval))
}

func (h *IntervalHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *IntervalHeap) Peek() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	return x
}

func copyEvent(e *eevent.Event) *eevent.Event {
	return &eevent.Event{
		ID:              e.ID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		From:            e.From,
		Till:            e.Till,
		CreatorID:       e.CreatorID,
		User:            e.User,
		ParticipantsIDs: e.ParticipantsIDs,
		Users:           e.Users,
		Invites:         e.Invites,
		Details:         e.Details,
		ScheduleRule:    e.ScheduleRule,
		EventSchedules:  e.EventSchedules,
		IsRepeat:        e.IsRepeat,
		IsPrivate:       e.IsPrivate,
	}
}
