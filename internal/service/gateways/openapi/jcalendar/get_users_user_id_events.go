package jcalendar

import (
	"net/http"
	"strconv"

	qrevent "jcalendar/internal/service/usecase/queries/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetUsersUserIdEvents(c echo.Context, userId string, params jcalendarsrv.GetUsersUserIdEventsParams) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(userId)
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
