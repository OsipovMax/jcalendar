package users

import (
	"context"
	"errors"
	"fmt"

	euser "jcalendar/internal/service/entity/user"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, u *euser.User) error {
	err := r.db.WithContext(ctx).Create(u).Error
	if err != nil {
		return fmt.Errorf("invalid creating user: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uint) (*euser.User, error) {
	u := &euser.User{}

	err := r.db.WithContext(ctx).First(u, id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("invalid getting user by id: %w", err)
	}

	return u, err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*euser.User, error) {
	u := &euser.User{}

	err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(u).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("invalid getting user by email: %w", err)
	}

	return u, err
}
