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

func (ch *CreateEventScheduleCommandHandler) Handle(ctx context.Context, command *CreateEventScheduleCommand) (uint, error) {
	//es := schedule.NewEventsSchedule(
	//	ctx,
	//	command.BeginOccurrence,
	//	command.EndOccurrence,
	//	command.EndingMode,
	//	command.Shift,
	//	command.SchedulerType,
	//	command.IntervalVal,
	//	command.IsRegular,
	//	command.EventID,
	//)
	//
	//_, err := ch.creator.CreateInvite(ctx, i)
	//if err != nil {
	//	return 0, err
	//}

	return 0, nil
}
