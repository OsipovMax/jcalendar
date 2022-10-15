package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	eevent "jcalendar/internal/service/entity/event"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/users"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func TestPostEvents(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/events"
		method = "POST"
	)

	table := []*struct {
		testSubTittle      string
		eventRequest       jcalendarsrv.EventRequest
		expectedStatusCode int
	}{
		{
			testSubTittle: "invalid request body",
			eventRequest: jcalendarsrv.EventRequest{
				Data: jcalendarsrv.InputEvent{
					From:      eventFromTimestamp.Format(time.RFC3339),
					Till:      eventTillTimestamp.Format(time.RFC3339),
					CreatorID: 0,
					Details:   "details",
					IsRepeat:  false,
					IsPrivate: true,
				},
			},
			expectedStatusCode: 400,
		},
		{
			testSubTittle: "successfully creating new event",
			eventRequest: jcalendarsrv.EventRequest{
				Data: jcalendarsrv.InputEvent{
					From:      eventFromTimestamp.Format(time.RFC3339),
					Till:      eventTillTimestamp.Format(time.RFC3339),
					CreatorID: 1,
					Details:   "details",
					IsRepeat:  false,
					IsPrivate: true,
				},
			},
		},
	}

	db, err := pkg.NewDB(ctx)
	require.NoError(t, err)

	var (
		urepo = users.NewRepository(db)
		erepo = events.NewRepository(db)
	)

	for _, row := range table {
		t.Run(row.testSubTittle, func(t *testing.T) {
			defer func() {
				require.NoError(t, clearTables(ctx, db))
			}()

			var application *app.Application
			application, err = app.NewApplication(ctx, db)
			require.NoError(t, err)

			creator := euser.NewUser(ctx, userFirstName, userLastName, userEmail, userhp, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, creator)
			require.NoError(t, err)

			var rawData []byte
			rawData, err = json.Marshal(row.eventRequest)
			require.NoError(t, err)

			req := httptest.NewRequest(method, path, bytes.NewBuffer(rawData))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).PostEvents(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatusCode, actualErr.Code)
			} else {
				actual := jcalendarsrv.CreatedEvent{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
				var actualEv *eevent.Event
				actualEv, err = erepo.GetEventByID(ctx, uint(actual.ID))
				require.NoError(t, err)

				var ft time.Time
				ft, err = time.Parse(time.RFC3339, row.eventRequest.Data.From)
				require.NoError(t, err)

				var tt time.Time
				tt, err = time.Parse(time.RFC3339, row.eventRequest.Data.Till)
				require.NoError(t, err)

				actualEv.User.HashedPassword = ""
				require.True(t, reflect.DeepEqual(
					eevent.Event{
						ID:        actualEv.ID,
						CreatedAt: actualEv.CreatedAt,
						UpdatedAt: actualEv.UpdatedAt,
						From:      ft,
						Till:      tt,
						CreatorID: uint(row.eventRequest.Data.CreatorID),
						User: &euser.User{
							ID:             actualEv.User.ID,
							CreatedAt:      actualEv.User.CreatedAt,
							UpdatedAt:      actualEv.User.UpdatedAt,
							FirstName:      creator.FirstName,
							LastName:       creator.LastName,
							Email:          creator.Email,
							TimeZoneOffset: creator.TimeZoneOffset,
						},
						Users:     []*euser.User{},
						Details:   row.eventRequest.Data.Details,
						IsRepeat:  row.eventRequest.Data.IsRepeat,
						IsPrivate: row.eventRequest.Data.IsPrivate,
					},
					*actualEv,
				))
			}
		})
	}
}
