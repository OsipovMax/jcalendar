package events

import (
	"context"
	"fmt"
	"time"

	eevent "jcalendar/internal/service/entity/event"
	einvite "jcalendar/internal/service/entity/invite"
	euser "jcalendar/internal/service/entity/user"

	"gorm.io/gorm"
)

const (
	ParticipantsAssociations = "Users"
	InvitesAssociations      = "Invites"
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
	/*
		----------------------------------------------------------------------------
	*/
	users := make([]*euser.User, 0, len(e.ParticipantsIDs))
	invites := make([]*einvite.Invite, 0, len(e.ParticipantsIDs))
	for i := range e.ParticipantsIDs {
		users = append(users, &euser.User{ID: e.ParticipantsIDs[i]})
		invites = append(invites, &einvite.Invite{UserID: e.ParticipantsIDs[i], IsAccepted: false})
	}

	e.Invites = invites

	/*
		----------------------------------------------------------------------------
	*/

	err := r.db.WithContext(ctx).Create(e).Error
	if err != nil {
		return 0, fmt.Errorf("invalid creating event: %w", err)
	}

	err = r.db.WithContext(ctx).Model(&e).
		Omit(fmt.Sprintf("%s.*", ParticipantsAssociations)).
		Association(ParticipantsAssociations).
		Append(users)
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

func (r *Repository) GetEventsInInterval(ctx context.Context, userID uint, from, till time.Time) ([]*eevent.Event, error) {
	es := make([]*eevent.Event, 0)
	err := r.db.WithContext(ctx).
		Preload(InvitesAssociations, "user_id = ? AND is_accepted = true", userID).
		Where("from >= ? AND till < ?", from, till).
		Find(es).
		Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("invalid getting events in interval by id: %w", err)
	}

	return es, err
}
