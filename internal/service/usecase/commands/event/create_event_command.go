package event

import (
	"context"
	"errors"
	"fmt"
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
	from, till string,
	creatorID uint,
	participantsIDs []int,
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

	tmp := make([]uint, len(participantsIDs), len(participantsIDs))
	for idx := range participantsIDs {
		tmp[idx] = uint(participantsIDs[idx])
	}

	ft, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return nil, fmt.Errorf("invalid converting from data: %w", err)
	}

	tt, err := time.Parse(time.RFC3339, till)
	if err != nil {
		return nil, fmt.Errorf("invalid converting till data: %w", err)
	}

	return &CreateEventCommand{
		From:            ft,
		Till:            tt,
		CreatorID:       creatorID,
		ParticipantsIDs: tmp,
		Details:         details,
		IsRepeat:        isRepeat,
		IsPrivate:       isPrivate,
	}, nil
}
