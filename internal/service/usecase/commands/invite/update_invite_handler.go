package invite

import (
	"context"
)

type Updater interface {
	UpdateInviteStatusByID(ctx context.Context, id uint, isAccepted bool) error
}

type UpdateInviteCommandHandler struct {
	updater Updater
}

func NewUpdateInviteCommandHandler(updater Updater) UpdateInviteCommandHandler {
	return UpdateInviteCommandHandler{updater: updater}
}

func (ch *UpdateInviteCommandHandler) Handle(ctx context.Context, command *UpdateInviteCommand) error {
	return ch.updater.UpdateInviteStatusByID(ctx, command.InviteID, command.IsAccepted)
}
