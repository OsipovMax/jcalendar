package jcalendar

import (
	"net/http"
	"strconv"

	qrevent "jcalendar/internal/service/usecase/queries/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetEventsId(c echo.Context, id string) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	query, err := qrevent.NewQuery(ctx, uint(iid))
	if err != nil {
		return err
	}

	e, err := s.application.Queries.GetEvent.Handle(ctx, query)
	if err != nil {
		return err
	}

	/*
		move ->
	*/
	var (
		uID     = c.Get("userID").(uint)
		Details = e.Details
	)

	if uID != e.CreatorID {
		Details = "Busy"
	}
	/**/

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.EventResponse{
			Data: &jcalendarsrv.OutputEvent{
				ID:       pcaster(int(e.ID)),
				CreateAt: pcaster(e.CreatedAt.String()),
				UpdateAt: pcaster(e.UpdatedAt.String()),
				From:     pcaster(e.From.String()),
				Till:     pcaster(e.Till.String()),
				Details:  &Details,
				Creator: &jcalendarsrv.OutputUser{
					ID:             pcaster(int(e.User.ID)),
					CreateAt:       pcaster(e.User.CreatedAt.String()),
					UpdateAt:       pcaster(e.User.UpdatedAt.String()),
					FirstName:      &e.User.FirstName,
					LastName:       &e.User.LastName,
					Email:          &e.User.Email,
					TimeZoneOffset: &e.User.TimeZoneOffset,
				},
				IsPrivate: &e.IsPrivate,
			},
		},
	)
}
