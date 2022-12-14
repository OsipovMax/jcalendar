package event

import (
	"context"
	"time"

	"jcalendar/internal/service/entity/invite"
	"jcalendar/internal/service/entity/schedule"
	"jcalendar/internal/service/entity/user"
)

type Event struct {
	ID              uint                      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
	From            time.Time                 `json:"from"`
	Till            time.Time                 `json:"till"`
	CreatorID       uint                      `json:"-"`
	User            *user.User                `json:"creator" gorm:"foreignKey:CreatorID;references:ID"`
	ParticipantsIDs []uint                    `json:"-" gorm:"-"`
	Users           []*user.User              `json:"users" gorm:"many2many:events_users"`
	Invites         []*invite.Invite          `json:"invites"`
	Details         string                    `json:"details"`
	ScheduleRule    *string                   `json:"schedule_rule"`
	EventSchedules  []*schedule.EventSchedule `json:"-"`
	IsPrivate       bool                      `json:"is_private"`
	IsRepeat        bool                      `json:"is_repeat"`
}

func NewEvent(
	_ context.Context,
	from, till time.Time,
	creatorID uint,
	participantsIDs []uint,
	details string,
	scheduleRule *string,
	eventSchedule []*schedule.EventSchedule,
	users []*user.User,
	invites []*invite.Invite,
	isPrivate, isRepeat bool,
) *Event {
	return &Event{
		From:            from,
		Till:            till,
		CreatorID:       creatorID,
		ParticipantsIDs: participantsIDs,
		Details:         details,
		EventSchedules:  eventSchedule,
		ScheduleRule:    scheduleRule,
		Users:           users,
		Invites:         invites,
		IsPrivate:       isPrivate,
		IsRepeat:        isRepeat,
	}
}
