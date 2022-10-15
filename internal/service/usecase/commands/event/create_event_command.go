package event

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type CreateEventCommand struct {
	From            time.Time
	Till            time.Time
	CreatorID       uint
	ParticipantsIDs []uint
	Details         string
	ScheduleRule    string
	IsPrivate       bool
	IsRepeat        bool
}

func NewCreateEventCommand(
	_ context.Context,
	from, till string,
	creatorID uint,
	participantsIDs []int,
	scheduleRule string,
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

	tmp := make([]uint, len(participantsIDs))
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
		ScheduleRule:    scheduleRule,
		IsRepeat:        isRepeat,
		IsPrivate:       isPrivate,
	}, nil
}
