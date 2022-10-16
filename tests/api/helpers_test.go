package api

import (
	"context"
	"time"

	"gorm.io/gorm"

	eevent "jcalendar/internal/service/entity/event"
	einvite "jcalendar/internal/service/entity/invite"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/repository/events"
	mrule "jcalendar/internal/service/usecase/managers/rule"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
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
	eventFromTimestamp = time.Date(2022, 1, 3, 12, 0, 0, 0, time.Local)
	eventTillTimestamp = time.Date(2022, 1, 3, 12, 30, 0, 0, time.Local)
	eventDetails       = "Details"

	userFirstName      = "Ivan"
	userLastName       = "Ivanov"
	userEmail          = "ivanov@mail.ru"
	userhp             = "5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5"
	userTimeZoneOffset = 10800

	inviteUserID  uint = 1
	inviteEventID uint = 1
)

func clearTables(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Exec(clearTablesQuery).Error
}

func getEventJSON(c, p *euser.User) jcalendarsrv.EventResponse {
	participants := []jcalendarsrv.OutputUser{
		{
			FirstName:      p.FirstName,
			LastName:       p.LastName,
			Email:          p.Email,
			TimeZoneOffset: p.TimeZoneOffset,
		},
	}

	return jcalendarsrv.EventResponse{
		Data: jcalendarsrv.OutputEvent{
			From:    eventFromTimestamp.String(),
			Till:    eventTillTimestamp.String(),
			Details: eventDetails,
			Creator: jcalendarsrv.OutputUser{
				FirstName:      c.FirstName,
				LastName:       c.LastName,
				Email:          c.Email,
				TimeZoneOffset: c.TimeZoneOffset,
			},
			Participants: &participants,
			IsPrivate:    false,
			IsRepeat:     false,
		},
	}
}

// nolint:funlen
func fillEvents(ctx context.Context, erepo *events.Repository, beginInterval, endInterval time.Time) error {
	/*

		E1(NR)_E2(R)___IB___E3(NR)____E4(R)____IE_____E5(NR)___E6(R)

	*/
	err := erepo.CreateEvent(ctx,
		eevent.NewEvent(
			ctx,
			beginInterval.AddDate(0, 0, -1),
			beginInterval.AddDate(0, 0, -1).Add(30*time.Minute),
			1,
			[]uint{2},
			"event 1",
			nil,
			nil,
			[]*euser.User{{ID: 2}},
			[]*einvite.Invite{{UserID: 2, IsAccepted: false}},
			false,
			false,
		),
	)
	if err != nil {
		return err
	}

	customScheduleRule := "SCHEDULER_MODE=CUSTOM;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=WEEKLY;CUSTOM_DAY_LIST=1,2"
	ruleManager := mrule.NewRuleManager(ctx)
	schedules, err := ruleManager.HandleRule(ctx, beginInterval.AddDate(0, 0, -1), customScheduleRule)
	if err != nil {
		return err
	}

	err = erepo.CreateEvent(ctx,
		eevent.NewEvent(
			ctx,
			beginInterval.AddDate(0, 0, -1),
			beginInterval.AddDate(0, 0, -1).Add(30*time.Minute),
			1,
			[]uint{2},
			"event 2",
			&customScheduleRule,
			schedules,
			nil,
			nil,
			false,
			true,
		),
	)
	if err != nil {
		return err
	}

	err = erepo.CreateEvent(ctx,
		eevent.NewEvent(
			ctx,
			beginInterval.Add(5*time.Hour),
			beginInterval.Add(5*time.Hour).Add(30*time.Minute),
			1,
			[]uint{2},
			"event 3",
			nil,
			nil,
			[]*euser.User{{ID: 2}},
			[]*einvite.Invite{{UserID: 2, IsAccepted: false}},
			false,
			false,
		),
	)
	if err != nil {
		return err
	}

	commonScheduleRule := "SCHEDULER_MODE=COMMON;ENDING_MODE=NONE;INTERVAL=1;IS_REGULAR=TRUE;SHIFT=DAILY"
	schedules, err = ruleManager.HandleRule(ctx, beginInterval.Add(5*time.Hour), commonScheduleRule)
	if err != nil {
		return err
	}

	err = erepo.CreateEvent(ctx,
		eevent.NewEvent(
			ctx,
			beginInterval.Add(5*time.Hour),
			beginInterval.Add(5*time.Hour).Add(30*time.Minute),
			1,
			[]uint{2},
			"event 4",
			&commonScheduleRule,
			schedules,
			nil,
			nil,
			false,
			true,
		),
	)
	if err != nil {
		return err
	}

	err = erepo.CreateEvent(ctx,
		eevent.NewEvent(
			ctx,
			endInterval.AddDate(0, 0, 1),
			endInterval.AddDate(0, 0, 1).Add(30*time.Minute),
			1,
			[]uint{2},
			"event 5",
			nil,
			nil,
			[]*euser.User{{ID: 2}},
			[]*einvite.Invite{{UserID: 2, IsAccepted: false}},
			false,
			false,
		),
	)
	if err != nil {
		return err
	}

	schedules, err = ruleManager.HandleRule(ctx, endInterval.AddDate(0, 0, 1), customScheduleRule)
	if err != nil {
		return err
	}
	err = erepo.CreateEvent(ctx,
		eevent.NewEvent(
			ctx,
			endInterval.AddDate(0, 0, 1),
			endInterval.AddDate(0, 0, 1).Add(30*time.Minute),
			1,
			[]uint{2},
			"event 6",
			&customScheduleRule,
			schedules,
			nil,
			nil,
			false,
			true,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
