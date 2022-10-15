package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"net/url"
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

func TestPostLogin(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/login"
		method = "POST"
	)

	table := []*struct {
		testSubTittle      string
		username, password string
		expectedStatusCode int
	}{
		{
			testSubTittle:      "not exist username",
			username:           "pavlov@mail.ru",
			password:           "12345",
			expectedStatusCode: 404,
		},
		{
			testSubTittle:      "wrong password",
			username:           userEmail,
			password:           "123123",
			expectedStatusCode: 401,
		},
		{
			testSubTittle: "correct password",
			username:      userEmail,
			password:      "12345",
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

			creator := euser.NewUser(ctx, userFirstName, userLastName, userEmail, userhp, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, creator)
			require.NoError(t, err)

			req := httptest.NewRequest(method, path, nil)
			req.Form = url.Values{}
			req.Form.Set("username", row.username)
			req.Form.Set("password", row.password)

			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).PostLogin(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatusCode, actualErr.Code)
			} else {
				actual := jcalendarsrv.TokenResponse{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))
				require.NotEqual(t, 0, len(actual.Data.Token))
			}
		})
	}
}
