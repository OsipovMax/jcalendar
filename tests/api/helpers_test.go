package api

import (
	"context"
	"time"

	euser "jcalendar/internal/service/entity/user"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"gorm.io/gorm"
)

const (
	clearTablesQuery = `
		DELETE FROM event_schedules;
		ALTER SEQUENCE event_schedules_id_seq RESTART WITH 1;
		DELETE FROM invites;
		ALTER SEQUENCE invites_id_seq RESTART WITH 1;
		DELETE FROM events_users;
		ALTER SEQUENCE events_users_id_seq RESTART WITH 1;
		DELETE FROM events;
		ALTER SEQUENCE events_id_seq RESTART WITH 1;
		DELETE FROM users;
		ALTER SEQUENCE users_id_seq RESTART WITH 1;
	`
)

var (
	eventFromTimestamp = time.Date(2022, 1, 1, 12, 0, 0, 0, time.Local)
	eventTillTimestamp = time.Date(2022, 1, 1, 12, 30, 0, 0, time.Local)
	eventIsPrivate     = false
	eventIsRepeat      = false
	eventDetails       = "Details"

	userFirstName      = "Ivan"
	userLastName       = "Ivanov"
	userEmail          = "ivanov@mail.ru"
	userHashedPassword = "5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5"
	userTimeZoneOffset = 10800

	inviteUserID  uint = 1
	inviteEventID uint = 1
)

func clearTables(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Exec(clearTablesQuery).Error
}

func pcaster[T comparable](val T) *T {
	return &val
}

func getEventJSON(c, p *euser.User) jcalendarsrv.EventResponse {
	participants := []jcalendarsrv.OutputUser{
		{
			FirstName:      &p.FirstName,
			LastName:       &p.LastName,
			Email:          &p.Email,
			TimeZoneOffset: &p.TimeZoneOffset,
		},
	}

	return jcalendarsrv.EventResponse{
		Data: &jcalendarsrv.OutputEvent{
			From:    pcaster(eventFromTimestamp.String()),
			Till:    pcaster(eventTillTimestamp.String()),
			Details: &eventDetails,
			Creator: &jcalendarsrv.OutputUser{
				FirstName:      &c.FirstName,
				LastName:       &c.LastName,
				Email:          &c.Email,
				TimeZoneOffset: &c.TimeZoneOffset,
			},
			Participants: &participants,
			IsPrivate:    &eventIsPrivate,
			IsRepeat:     &eventIsRepeat,
		},
	}
}
