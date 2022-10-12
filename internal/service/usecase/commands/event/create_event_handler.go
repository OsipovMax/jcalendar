package event

import (
	"context"

	eevent "jcalendar/internal/service/entity/event"
)

type Creator interface {
	CreateEvent(ctx context.Context, e *eevent.Event) (uint, error)
}

type CreateEventCommandHandler struct {
	creator Creator
}

func NewCreateEventCommandHandler(creator Creator) CreateEventCommandHandler {
	return CreateEventCommandHandler{creator: creator}
}

func (ch *CreateEventCommandHandler) Handle(ctx context.Context, command *CreateEventCommand) (uint, error) {
	e := eevent.NewEvent(
		ctx,
		command.From,
		command.Till,
		command.CreatorID,
		command.ParticipantsIDs,
		command.Details,
		command.IsPrivate,
		command.IsRepeat,
	)

	_, err := ch.creator.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}

	//for _, participantID := range e.ParticipantsIDs {
	//	var icmd *invite.CreateInviteCommand
	//	icmd, err = invite.NewCreateInviteCommand(ctx, participantID, eID, false)
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = s.application.Commands.CreateInvite.Handle(ctx, icmd)
	//	if err != nil {
	//		return err
	//	}
	//}

	return e.ID, nil
}
