package jcalendar

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"jcalendar/internal/service/app"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestGetEventsId(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/events/:id"
		method = "GET"
		//headers = map[string]string{"Content-Type": "application/json; charset=UTF-8"}
	)

	table := []*struct {
		TestSubTittle      string
		ID                 string
		ExpectedStatusCode int
	}{
		{
			TestSubTittle:      "invalid event id format",
			ID:                 "abc",
			ExpectedStatusCode: 500,
		},
	}

	for _, row := range table {
		t.Run(row.TestSubTittle, func(t *testing.T) {
			application, err := app.NewApplication(ctx)
			require.Error(t, err)

			req := httptest.NewRequest(method, path, nil)
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("asd")

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: NewServer(application)}).GetEventsId(c)
			require.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	}
}
