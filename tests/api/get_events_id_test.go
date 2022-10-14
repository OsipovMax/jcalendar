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
				actual := jcalendarsrv.EventResponse{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))

				participants := []jcalendarsrv.OutputUser{
					{
						ID:             (*actual.Data.Participants)[0].ID,
						CreateAt:       (*actual.Data.Participants)[0].CreateAt,
						UpdateAt:       (*actual.Data.Participants)[0].UpdateAt,
						FirstName:      &participant.FirstName,
						LastName:       &participant.LastName,
						Email:          &participant.Email,
						TimeZoneOffset: &participant.TimeZoneOffset,
					},
				}

				require.True(t, reflect.DeepEqual(
					jcalendarsrv.EventResponse{
						Data: &jcalendarsrv.OutputEvent{
							ID:       pcaster(1),
							CreateAt: actual.Data.CreateAt,
							UpdateAt: actual.Data.UpdateAt,
							From:     pcaster(eventFromTimestamp.String()),
							Till:     pcaster(eventTillTimestamp.String()),
							Details:  &eventDetails,
							Creator: &jcalendarsrv.OutputUser{
								ID:             actual.Data.Creator.ID,
								CreateAt:       actual.Data.Creator.CreateAt,
								UpdateAt:       actual.Data.Creator.UpdateAt,
								FirstName:      &creator.FirstName,
								LastName:       &creator.LastName,
								Email:          &creator.Email,
								TimeZoneOffset: &creator.TimeZoneOffset,
							},
							Participants: &participants,
							IsPrivate:    &eventIsPrivate,
							IsRepeat:     &eventIsRepeat,
						},
					},
					actual,
				))
			}
		})
	}
}
