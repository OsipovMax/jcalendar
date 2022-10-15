package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	eevent "jcalendar/internal/service/entity/event"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/users"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetEventsId(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/events/:id"
		method = "GET"
	)

	table := []*struct {
		testSubTittle      string
		id                 string
		expectedStatusCode int
	}{
		{
			testSubTittle:      "invalid event id format",
			id:                 "abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testSubTittle:      "invalid event id value",
			id:                 "0",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testSubTittle:      "not existing event",
			id:                 "5",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			testSubTittle:      "successfully getting event",
			id:                 "1",
			expectedStatusCode: http.StatusOK,
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

			creator := euser.NewUser(ctx, userFirstName, userLastName, userEmail, userHashedPassword, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, creator)
			require.NoError(t, err)

			participant := euser.NewUser(ctx, userFirstName, "sidorov", "sidorov@mail.ru", userHashedPassword, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, participant)
			require.NoError(t, err)

			_, err = erepo.CreateEvent(ctx,
				eevent.NewEvent(
					ctx,
					eventFromTimestamp,
					eventTillTimestamp,
					1,
					[]uint{2},
					eventDetails,
					nil,
					false,
					false,
				),
			)

			req := httptest.NewRequest(method, path, nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(row.id)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).GetEventsId(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatusCode, actualErr.Code)
			} else {
				actualEvent := jcalendarsrv.EventResponse{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actualEvent))

				expectedEvent := getEventJSON(creator, participant)

				expectedEvent.Data.ID = pcaster(1)
				expectedEvent.Data.CreateAt = actualEvent.Data.CreateAt
				expectedEvent.Data.UpdateAt = actualEvent.Data.UpdateAt

				expectedEvent.Data.Creator.ID = actualEvent.Data.Creator.ID
				expectedEvent.Data.Creator.CreateAt = actualEvent.Data.Creator.CreateAt
				expectedEvent.Data.Creator.UpdateAt = actualEvent.Data.Creator.UpdateAt

				(*expectedEvent.Data.Participants)[0].ID = (*actualEvent.Data.Participants)[0].ID
				(*expectedEvent.Data.Participants)[0].CreateAt = (*actualEvent.Data.Participants)[0].CreateAt
				(*expectedEvent.Data.Participants)[0].UpdateAt = (*actualEvent.Data.Participants)[0].UpdateAt

				require.True(t, reflect.DeepEqual(expectedEvent, actualEvent))
			}
		})
	}
}
