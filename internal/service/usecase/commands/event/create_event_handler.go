package event

import (
	"context"

	eevent "jcalendar/internal/service/entity/event"
	"jcalendar/internal/service/entity/schedule"
	"jcalendar/internal/service/usecase/manager"
)

type Creator interface {
	CreateEvent(ctx context.Context, e *eevent.Event) (uint, error)
}

// TODO: pass handler only
type CreateEventCommandHandler struct {
	creator      Creator
	eventManager *manager.EventManager
}

func NewCreateEventCommandHandler(creator Creator, eventManager *manager.EventManager) CreateEventCommandHandler {
	return CreateEventCommandHandler{creator: creator, eventManager: eventManager}
}

func (ch *CreateEventCommandHandler) Handle(ctx context.Context, command *CreateEventCommand) (uint, error) {
	var schedules []*schedule.EventSchedule
	if command.IsRepeat {
		var err error
		schedules, err = ch.eventManager.HandleRule(ctx, command.From, command.ScheduleRule)
		if err != nil {
			return 0, err
		}
	}

	e := eevent.NewEvent(
		ctx,
		command.From,
		command.Till,
		command.CreatorID,
		command.ParticipantsIDs,
		command.Details,
		schedules,
		command.IsPrivate,
		command.IsRepeat,
	)

	_, err := ch.creator.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}

	return e.ID, nil
}
