package event

import (
	"context"
	"errors"
)

type GetInviteQuery struct {
	InviteID uint
}

func NewGetInviteQuery(_ context.Context, id uint) (*GetInviteQuery, error) {
	if id <= 0 {
		return nil, errors.New("non-positive id value")
	}

	return &GetInviteQuery{
		InviteID: id,
	}, nil
}
