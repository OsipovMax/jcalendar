package invite

import (
	"context"

	einvite "jcalendar/internal/service/entity/invite"
)

type Creator interface {
	CreateInvite(ctx context.Context, i *einvite.Invite) (uint, error)
}

type CreateInviteCommandHandler struct {
	creator Creator
}

func NewCreateInviteCommandHandler(creator Creator) CreateInviteCommandHandler {
	return CreateInviteCommandHandler{creator: creator}
}

func (ch *CreateInviteCommandHandler) Handle(ctx context.Context, command *CreateInviteCommand) (uint, error) {
	i := einvite.NewInvite(ctx, command.UserID, command.EventID, command.IsAccepted)

	_, err := ch.creator.CreateInvite(ctx, i)
	if err != nil {
		return 0, err
	}

	return i.ID, nil
}
