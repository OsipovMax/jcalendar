package user

import (
	"context"

	euser "jcalendar/internal/service/entity/user"
)

type Getter interface {
	GetUserByEmail(ctx context.Context, email string) (*euser.User, error)
}

type GetUserByEmailQueryHandler struct {
	getter Getter
}

func NewGetUserByEmailQueryHandler(getter Getter) GetUserByEmailQueryHandler {
	return GetUserByEmailQueryHandler{getter: getter}
}

func (ch *GetUserByEmailQueryHandler) Handle(ctx context.Context, query *GetUserByEmailQuery) (*euser.User, error) {
	e, err := ch.getter.GetUserByEmail(ctx, query.UserEmail)
	if err != nil {
		return nil, err
	}

	return e, nil
}
