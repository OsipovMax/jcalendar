package invites

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	einvite "jcalendar/internal/service/entity/invite"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateInvite(ctx context.Context, i *einvite.Invite) error {
	err := r.db.WithContext(ctx).Create(i).Error

	if err != nil {
		return fmt.Errorf("invalid creating invite: %w", err)
	}

	return nil
}

func (r *Repository) GetInviteByID(ctx context.Context, id uint) (*einvite.Invite, error) {
	inv := &einvite.Invite{}
	err := r.db.WithContext(ctx).First(inv, id).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("invalid getting invite by id: %w", err)
	}

	return inv, err
}

func (r *Repository) UpdateInviteStatusByID(ctx context.Context, id uint, isAccepted bool) error {
	err := r.db.WithContext(ctx).
		Model(&einvite.Invite{ID: id}).
		Update("is_accepted", isAccepted).
		Error

	if err != nil {
		return fmt.Errorf("invalid updating invite status: %w", err)
	}

	return nil
}
