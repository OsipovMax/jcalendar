package event

import (
	"context"

	eevent "jcalendar/internal/service/entity/event"
	"jcalendar/internal/service/usecase/parser"
)

type Creator interface {
	CreateEvent(ctx context.Context, e *eevent.Event) (uint, error)
}

type CreateEventCommandHandler struct {
	creator Creator
	sparser *parser.ScheduleParser
}

func NewCreateEventCommandHandler(creator Creator, sparser *parser.ScheduleParser) CreateEventCommandHandler {
	return CreateEventCommandHandler{creator: creator, sparser: sparser}
}

func (ch *CreateEventCommandHandler) Handle(ctx context.Context, command *CreateEventCommand) (uint, error) {
	rules, err := ch.sparser.HandleRule(ctx, command.From, command.ScheduleRule)
	if err != nil {
		return 0, err
	}

	e := eevent.NewEvent(
		ctx,
		command.From,
		command.Till,
		command.CreatorID,
		command.ParticipantsIDs,
		command.Details,
		rules,
		command.IsPrivate,
		command.IsRepeat,
	)

	_, err = ch.creator.CreateEvent(ctx, e)
	if err != nil {
		return 0, err
	}

	return e.ID, nil
}
