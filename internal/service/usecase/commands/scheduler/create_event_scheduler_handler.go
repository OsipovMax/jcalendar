package scheduler

import (
	"context"
)

type Creator interface {
	CreateSchedule(ctx context.Context) (uint, error)
}

type CreateEventScheduleCommandHandler struct {
	creator Creator
}

func NewCreateEventScheduleCommandHandler(creator Creator) CreateEventScheduleCommandHandler {
	return CreateEventScheduleCommandHandler{creator: creator}
}

func (ch *CreateEventScheduleCommandHandler) Handle(ctx context.Context, command *CreateEventScheduleCommandHandler) (uint, error) {
	return 0, nil
}
