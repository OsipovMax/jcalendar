package event

import (
	"context"
	"errors"
	"time"

	euser "jcalendar/internal/service/entity/user"
)

type CreateEventCommand struct {
	From            time.Time
	Till            time.Time
	CreatorID       uint
	Creator         *euser.User
	ParticipantsIDs []uint
	Details         string
	IsPrivate       bool
	IsRepeat        bool
}

func NewCreateEventCommand(
	_ context.Context,
	from, till time.Time,
	creatorID uint,
	participantsIDs []int,
	details string,
	isPrivate, isRepeat bool,
) (*CreateEventCommand, error) {
	if from.IsZero() {
		return nil, errors.New("missing from value")
	}

	if till.IsZero() {
		return nil, errors.New("missing till value")
	}

	if creatorID <= 0 {
		return nil, errors.New("non-positive creatorID value")
	}

	tmp := make([]uint, len(participantsIDs), len(participantsIDs))
	for idx := range participantsIDs {
		tmp[idx] = uint(participantsIDs[idx])
	}

	return &CreateEventCommand{
		From:            from,
		Till:            till,
		CreatorID:       creatorID,
		ParticipantsIDs: tmp,
		Details:         details,
		IsRepeat:        isRepeat,
		IsPrivate:       isPrivate,
	}, nil
}
