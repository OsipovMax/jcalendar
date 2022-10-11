package invite

import (
	"context"
	"errors"
)

type UpdateInviteCommand struct {
	InviteID   uint
	IsAccepted bool
}

func NewUpdateInviteCommand(_ context.Context, inviteID uint, isAccepted bool) (*UpdateInviteCommand, error) {
	if inviteID <= 0 {
		return nil, errors.New("non-positive userID value")
	}

	return &UpdateInviteCommand{
		InviteID:   inviteID,
		IsAccepted: isAccepted,
	}, nil
}
