package event

import (
	"errors"
)

type GetInviteQuery struct {
	InviteID uint
}

func NewInviteQuery(id uint) (*GetInviteQuery, error) {
	if id <= 0 {
		return nil, errors.New("non-positive id value")
	}

	return &GetInviteQuery{
		InviteID: id,
	}, nil
}
