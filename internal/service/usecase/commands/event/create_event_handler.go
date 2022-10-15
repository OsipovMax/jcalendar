package event

import (
	"context"

	eevent "jcalendar/internal/service/entity/event"
	einvite "jcalendar/internal/service/entity/invite"
	"jcalendar/internal/service/entity/schedule"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/usecase/manager"
)

type Creator interface {
	CreateEvent(ctx context.Context, e *eevent.Event) error
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

	users := make([]*euser.User, 0, len(command.ParticipantsIDs))
	invites := make([]*einvite.Invite, 0, len(command.ParticipantsIDs))
	for idx := range command.ParticipantsIDs {
		users = append(users, &euser.User{ID: command.ParticipantsIDs[idx]})
		invites = append(invites, &einvite.Invite{UserID: command.ParticipantsIDs[idx], IsAccepted: false})
	}

	e := eevent.NewEvent(
		ctx,
		command.From,
		command.Till,
		command.CreatorID,
		command.ParticipantsIDs,
		command.Details,
		schedules,
		users,
		invites,
		command.IsPrivate,
		command.IsRepeat,
	)

	err := ch.creator.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}

	return e.ID, nil
}
