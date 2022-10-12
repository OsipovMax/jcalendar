package invite

import (
	"context"
	"time"

	"jcalendar/internal/service/entity/user"
)

type Invite struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserID    uint       `json:"-"`
	User      *user.User `json:"user"`
	EventID   uint       `json:"-"`
	//Event      *event.Event `json:"event"`
	IsAccepted bool `json:"is_accepted"`
}

func NewInvite(_ context.Context, userID, eventID uint, isAccepted bool) *Invite {
	return &Invite{
		UserID:     userID,
		EventID:    eventID,
		IsAccepted: isAccepted,
	}
}
