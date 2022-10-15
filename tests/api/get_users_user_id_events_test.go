package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	eevent "jcalendar/internal/service/entity/event"
	einvite "jcalendar/internal/service/entity/invite"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/users"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

// nolint:funlen
func TestGetUsersUserIdEvents(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/users/:user_id/events"
		method = "GET"
	)

	table := []*struct {
		testSubTittle      string
		userID             string
		from, till         string
		expectedStatusCode int
		expectedEventsLen  int
	}{
		{
			testSubTittle:      "invalid user id format",
			userID:             "abc",
			from:               eventFromTimestamp.Format(time.RFC3339),
			till:               eventTillTimestamp.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
			expectedEventsLen:  0,
		},
		{
			testSubTittle:      "invalid user id value",
			userID:             "0",
			from:               eventFromTimestamp.Format(time.RFC3339),
			till:               eventTillTimestamp.Format(time.RFC3339),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testSubTittle:      "not existing user",
			userID:             "5",
			from:               eventFromTimestamp.Format(time.RFC3339),
			till:               eventTillTimestamp.Format(time.RFC3339),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			testSubTittle:      "successfully getting all events in interval",
			userID:             "1",
			from:               eventFromTimestamp.Format(time.RFC3339),
			till:               eventTillTimestamp.AddDate(0, 0, 10).Format(time.RFC3339),
			expectedStatusCode: http.StatusOK,
			expectedEventsLen:  2,
		},
		{
			testSubTittle:      "successfully getting single events in interval",
			userID:             "1",
			from:               eventFromTimestamp.Format(time.RFC3339),
			till:               eventTillTimestamp.AddDate(0, 0, 1).Format(time.RFC3339),
			expectedStatusCode: http.StatusOK,
			expectedEventsLen:  1,
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

			otherCreator := euser.NewUser(ctx, userFirstName, "sidorov", "sidorov@mail.ru", userhp, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, otherCreator)
			require.NoError(t, err)

			require.NoError(t,
				erepo.CreateEvent(ctx,
					eevent.NewEvent(
						ctx,
						eventFromTimestamp,
						eventTillTimestamp,
						1,
						[]uint{2},
						eventDetails,
						nil,
						nil,
						[]*euser.User{{ID: 2}},
						[]*einvite.Invite{{UserID: 2, IsAccepted: false}},
						false,
						false,
					),
				),
			)

			require.NoError(t,
				erepo.CreateEvent(ctx,
					eevent.NewEvent(
						ctx,
						eventFromTimestamp.AddDate(0, 0, 3),
						eventTillTimestamp.AddDate(0, 0, 3),
						1,
						[]uint{2},
						eventDetails,
						nil,
						nil,
						[]*euser.User{{ID: 2}},
						[]*einvite.Invite{{UserID: 2, IsAccepted: false}},
						false,
						false,
					),
				),
			)

			require.NoError(t,
				erepo.CreateEvent(ctx,
					eevent.NewEvent(
						ctx,
						eventFromTimestamp.AddDate(0, 0, 3),
						eventTillTimestamp.AddDate(0, 0, 3),
						2,
						[]uint{},
						eventDetails,
						nil,
						nil,
						nil,
						nil,
						false,
						false,
					),
				),
			)

			req := httptest.NewRequest(method, path, nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetParamNames("user_id")
			c.SetParamValues(row.userID)
			c.QueryParams().Set("from", row.from)
			c.QueryParams().Set("till", row.till)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).GetUsersUserIdEvents(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatusCode, actualErr.Code)
			} else {
				actualEvents := jcalendarsrv.EventsResponse{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actualEvents))
				require.Equal(t, row.expectedEventsLen, len(actualEvents.Data))
			}
		})
	}
}
