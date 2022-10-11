package event

import (
	"context"
	"errors"

	euser "jcalendar/internal/service/entity/user"
)

type CreateEventCommand struct {
	From            string
	Till            string
	CreatorID       uint
	Creator         *euser.User
	ParticipantsIDs []uint
	Details         string
	IsPrivate       bool
	IsRepeat        bool
}

func NewCreateEventCommand(
	_ context.Context,
	from, till string,
	creatorID uint,
	participantsIDs []uint,
	details string,
	isPrivate, isRepeat bool,
) (*CreateEventCommand, error) {
	if from == "" {
		return nil, errors.New("missing from value")
	}

	if till == "" {
		return nil, errors.New("missing till value")
	}

	if creatorID <= 0 {
		return nil, errors.New("non-positive creatorID value")
	}

	return &CreateEventCommand{
		From:            from,
		Till:            till,
		CreatorID:       creatorID,
		ParticipantsIDs: participantsIDs,
		Details:         details,
		IsRepeat:        isRepeat,
		IsPrivate:       isPrivate,
	}, nil
}
