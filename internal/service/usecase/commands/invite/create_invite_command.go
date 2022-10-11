package invite

import (
	"context"
	"errors"
)

type CreateInviteCommand struct {
	UserID     uint
	EventID    uint
	IsAccepted bool
}

func NewCreateInviteCommand(_ context.Context, userID, eventID uint, isAccepted bool) (*CreateInviteCommand, error) {
	if userID <= 0 {
		return nil, errors.New("non-positive userID value")
	}

	if eventID <= 0 {
		return nil, errors.New("non-positive eventID value")
	}

	return &CreateInviteCommand{
		UserID:     userID,
		EventID:    eventID,
		IsAccepted: isAccepted,
	}, nil
}
