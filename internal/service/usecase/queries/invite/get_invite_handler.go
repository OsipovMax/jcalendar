package event

import (
	"context"

	einvite "jcalendar/internal/service/entity/invite"
)

type Getter interface {
	GetInviteByID(ctx context.Context, id uint) (*einvite.Invite, error)
}

type GetInviteQueryHandler struct {
	getter Getter
}

func NewGetInviteQueryHandler(getter Getter) GetInviteQueryHandler {
	return GetInviteQueryHandler{getter: getter}
}

func (ch *GetInviteQueryHandler) Handle(ctx context.Context, query *GetInviteQuery) (*einvite.Invite, error) {
	e, err := ch.getter.GetInviteByID(ctx, query.InviteID)
	if err != nil {
		return nil, err
	}

	return e, nil
}
