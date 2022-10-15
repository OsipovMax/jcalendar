package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestGetWindows(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/windows"
		method = "GET"
	)

	table := []*struct {
		testSubTittle              string
		userIDs                    []int
		winSize                    string
		expectedStatusCode         int
		expectedFrom, expectedTill string
	}{
		{
			testSubTittle:      "invalid win size param",
			userIDs:            []int{1, 2},
			winSize:            "abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testSubTittle:      "not existing users",
			userIDs:            []int{5, 7},
			winSize:            "15m",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			testSubTittle:      "successfully getting all events in interval",
			userIDs:            []int{1, 2},
			winSize:            "15m",
			expectedStatusCode: http.StatusOK,
			expectedFrom:       eventTillTimestamp.Format(time.RFC3339),
			expectedTill:       eventTillTimestamp.Add(15 * time.Minute).Format(time.RFC3339),
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

			ft := time.Now()
			tt := ft.Add(30 * time.Minute)
			require.NoError(t,
				erepo.CreateEvent(ctx,
					eevent.NewEvent(
						ctx,
						ft,
						tt,
						1,
						[]uint{},
						eventDetails,
						nil,
						nil,
						nil,
						false,
						false,
					),
				),
			)

			require.NoError(t,
				erepo.CreateEvent(ctx,
					eevent.NewEvent(
						ctx,
						ft.Add(time.Hour),
						tt.Add(time.Hour),
						2,
						[]uint{},
						eventDetails,
						nil,
						nil,
						nil,
						false,
						false,
					),
				),
			)

			req := httptest.NewRequest(
				method,
				fmt.Sprintf(
					"%s?user_ids[]=%d&user_ids[]=%d&win_size=%s",
					path, row.userIDs[0], row.userIDs[1], row.winSize,
				),
				nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).GetWindows(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatusCode, actualErr.Code)
			} else {
				actual := jcalendarsrv.FreeWindowResponse{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
				require.Equal(t, tt.Format(time.RFC3339), actual.Data.From)
				require.Equal(t, tt.Add(15*time.Minute).Format(time.RFC3339), actual.Data.Till)
			}
		})
	}
}
