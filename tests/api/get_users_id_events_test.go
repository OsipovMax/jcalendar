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
		id                 string
		from, till         time.Time
		expectedStatusCode int
		expectedEventsLen  int
	}{
		{
			testSubTittle:      "invalid user id format",
			id:                 "abc",
			from:               eventFromTimestamp,
			till:               eventTillTimestamp,
			expectedStatusCode: http.StatusBadRequest,
			expectedEventsLen:  0,
		},
		{
			testSubTittle:      "invalid user id value",
			id:                 "0",
			from:               eventFromTimestamp,
			till:               eventTillTimestamp,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testSubTittle:      "not existing user",
			id:                 "5",
			from:               eventFromTimestamp,
			till:               eventTillTimestamp,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			testSubTittle:      "successfully getting all events in interval(creator)",
			id:                 "1",
			from:               eventFromTimestamp,
			till:               eventFromTimestamp.AddDate(0, 0, 7),
			expectedStatusCode: http.StatusOK,
			expectedEventsLen:  10,
		},
		{
			testSubTittle:      "successfully getting single events in interval(participant)",
			id:                 "2",
			from:               eventFromTimestamp,
			till:               eventFromTimestamp.AddDate(0, 0, 7),
			expectedStatusCode: http.StatusOK,
			expectedEventsLen:  8,
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

			require.NoError(t, fillEvents(ctx, erepo, row.from, row.till))

			req := httptest.NewRequest(method, path, nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(row.id)
			c.QueryParams().Set("from", row.from.Format(time.RFC3339))
			c.QueryParams().Set("till", row.till.Format(time.RFC3339))
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).GetUsersIdEvents(c)
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
