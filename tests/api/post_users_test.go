package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	"jcalendar/internal/service/repository/users"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func TestPostUsers(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/users"
		method = "POST"
	)

	table := []*struct {
		testSubTittle      string
		userRequest        jcalendarsrv.UserRequest
		expectedStatusCode int
	}{
		{
			testSubTittle: "invalid request body",
			userRequest: jcalendarsrv.UserRequest{
				Data: jcalendarsrv.InputUser{
					FirstName:      userFirstName,
					LastName:       userLastName,
					Email:          userEmail,
					Password:       "123123",
					TimeZoneOffset: -1999,
				},
			},
			expectedStatusCode: 400,
		},
		{
			testSubTittle: "successfully creating new event",
			userRequest: jcalendarsrv.UserRequest{
				Data: jcalendarsrv.InputUser{
					FirstName:      userFirstName,
					LastName:       userLastName,
					Email:          userEmail,
					Password:       "123123",
					TimeZoneOffset: userTimeZoneOffset,
				},
			},
		},
	}

	db, err := pkg.NewDB(ctx)
	require.NoError(t, err)

	var (
		urepo = users.NewRepository(db)
	)

	for _, row := range table {
		t.Run(row.testSubTittle, func(t *testing.T) {
			defer func() {
				require.NoError(t, clearTables(ctx, db))
			}()

			var application *app.Application
			application, err = app.NewApplication(ctx, db)
			require.NoError(t, err)

			var rawData []byte
			rawData, err = json.Marshal(row.userRequest)
			require.NoError(t, err)

			req := httptest.NewRequest(method, path, bytes.NewBuffer(rawData))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).PostUsers(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatusCode, actualErr.Code)
			} else {
				actual := jcalendarsrv.CreatedUser{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
				var actualUser *euser.User
				actualUser, err = urepo.GetUserByID(ctx, uint(actual.ID))
				require.NoError(t, err)

				actualUser.HashedPassword = ""
				require.True(t, reflect.DeepEqual(
					euser.User{
						ID:             actualUser.ID,
						CreatedAt:      actualUser.CreatedAt,
						UpdatedAt:      actualUser.UpdatedAt,
						FirstName:      row.userRequest.Data.FirstName,
						LastName:       row.userRequest.Data.LastName,
						Email:          row.userRequest.Data.Email,
						TimeZoneOffset: row.userRequest.Data.TimeZoneOffset,
					},
					*actualUser,
				))
			}
		})
	}
}
