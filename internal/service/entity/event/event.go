package event

import (
	"context"
	"time"

	"jcalendar/internal/service/entity/user"
)

type Event struct {
	ID              uint         `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	From            string       `json:"from"`
	Till            string       `json:"till"`
	CreatorID       uint         `json:"-"`
	Creator         *user.User   `json:"creator"`
	ParticipantsIDs []uint       `json:"-" gorm:"-"`
	Participants    []*user.User `json:"brands" gorm:"many2many:events_users"`
	Details         string       `json:"details"`
	IsPrivate       bool         `json:"is_private"`
	IsRepeat        bool         `json:"is_repeat"`
}

func NewEvent(
	_ context.Context,
	from, till string,
	creatorID uint,
	participantsIDs []uint,
	details string,
	isPrivate, isRepeat bool,
) *Event {
	return &Event{
		From:            from,
		Till:            till,
		CreatorID:       creatorID,
		ParticipantsIDs: participantsIDs,
		Details:         details,
		IsPrivate:       isPrivate,
		IsRepeat:        isRepeat,
	}
}
