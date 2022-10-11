package user

import (
	"context"

	euser "jcalendar/internal/service/entity/user"
)

type Creator interface {
	CreateUser(ctx context.Context, user *euser.User) error
}

type CreateUserCommandHandler struct {
	creator Creator
}

func NewCommandHandler(creator Creator) CreateUserCommandHandler {
	return CreateUserCommandHandler{creator: creator}
}

func (ch *CreateUserCommandHandler) Handle(ctx context.Context, command *CreateUserCommand) (uint, error) {
	u := euser.NewUser(
		ctx,
		command.FirstName,
		command.LastName,
		command.Email,
		command.TimeZoneOffset,
	)

	err := ch.creator.CreateUser(ctx, u)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}
