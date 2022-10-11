package jcalendar

import (
	"net/http"
	"strconv"

	qrevent "jcalendar/internal/service/usecase/queries/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetEventsUserId(c echo.Context, id string, params jcalendarsrv.GetEventsUserIdParams) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	_, err = qrevent.NewQuery(ctx, uint(iid))
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.EventsResponse{},
	)
}
