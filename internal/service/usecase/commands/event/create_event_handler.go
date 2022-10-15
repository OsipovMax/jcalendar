package event

import (
	"context"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	einvite "jcalendar/internal/service/entity/invite"
	eschedule "jcalendar/internal/service/entity/schedule"
	euser "jcalendar/internal/service/entity/user"
)

type Creator interface {
	CreateEvent(ctx context.Context, e *eevent.Event) error
}

type RuleHandler interface {
	HandleRule(ctx context.Context, eventFrom time.Time, eventScheduleRule string) ([]*eschedule.EventSchedule, error)
}

type CreateEventCommandHandler struct {
	creator Creator
	handler RuleHandler
}

func NewCreateEventCommandHandler(creator Creator, handler RuleHandler) CreateEventCommandHandler {
	return CreateEventCommandHandler{creator: creator, handler: handler}
}

func (ch *CreateEventCommandHandler) Handle(ctx context.Context, command *CreateEventCommand) (uint, error) {
	var schedules []*eschedule.EventSchedule
	if command.IsRepeat && command.ScheduleRule != nil {
		var err error
		schedules, err = ch.handler.HandleRule(ctx, command.From, *command.ScheduleRule)
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
		command.ScheduleRule,
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
