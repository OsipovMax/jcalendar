package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

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

			creator := euser.NewUser(ctx, userFirstName, userLastName, userEmail, userhp, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, creator)
			require.NoError(t, err)

			participant := euser.NewUser(ctx, userFirstName, "sidorov", "sidorov@mail.ru", userhp, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, participant)
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
						[]*euser.User{{ID: 2}},
						[]*einvite.Invite{{UserID: 2, IsAccepted: false}},
						false,
						false,
					),
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

				expectedEvent.Data.ID = 1
				expectedEvent.Data.CreatedAt = actualEvent.Data.CreatedAt
				expectedEvent.Data.UpdatedAt = actualEvent.Data.UpdatedAt

				expectedEvent.Data.Creator.ID = actualEvent.Data.Creator.ID
				expectedEvent.Data.Creator.CreatedAt = actualEvent.Data.Creator.CreatedAt
				expectedEvent.Data.Creator.UpdatedAt = actualEvent.Data.Creator.UpdatedAt

				(*expectedEvent.Data.Participants)[0].ID = (*actualEvent.Data.Participants)[0].ID
				(*expectedEvent.Data.Participants)[0].CreatedAt = (*actualEvent.Data.Participants)[0].CreatedAt
				(*expectedEvent.Data.Participants)[0].UpdatedAt = (*actualEvent.Data.Participants)[0].UpdatedAt

				require.True(t, reflect.DeepEqual(expectedEvent, actualEvent))
			}
		})
	}
}
