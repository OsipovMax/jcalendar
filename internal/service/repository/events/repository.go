package events

import (
	"context"
	"fmt"

	eevent "jcalendar/internal/service/entity/event"
	euser "jcalendar/internal/service/entity/user"

	"gorm.io/gorm"
)

const (
	ParticipantsAssociations = "Users"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateEvent(ctx context.Context, e *eevent.Event) (uint, error) {
	participants := make([]*euser.User, 0, len(e.ParticipantsIDs))
	for i := range e.ParticipantsIDs {
		participants = append(participants, &euser.User{ID: uint(e.ParticipantsIDs[i])})
	}

	err := r.db.WithContext(ctx).Create(e).Error
	if err != nil {
		return 0, fmt.Errorf("invalid creating event: %w", err)
	}

	err = r.db.WithContext(ctx).Model(&e).
		Omit(fmt.Sprintf("%s.*", ParticipantsAssociations)).
		Association(ParticipantsAssociations).
		Append(participants)
	if err != nil {
		return 0, fmt.Errorf("invalid appending associated participants: %w", err)
	}

	return e.ID, nil
}

func (r *Repository) GetEventByID(ctx context.Context, id uint) (*eevent.Event, error) {
	e := &eevent.Event{}
	err := r.db.WithContext(ctx).
		Preload(ParticipantsAssociations).
		First(e, id).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("invalid getting event by id: %w", err)
	}

	return e, err
}
