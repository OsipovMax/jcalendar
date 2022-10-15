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
	return ch.getter.GetInviteByID(ctx, query.InviteID)
}
