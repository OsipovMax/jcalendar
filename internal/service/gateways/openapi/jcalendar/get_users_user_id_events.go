package jcalendar

import (
	"net/http"
	"strconv"

	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetUsersUserIdEvents(c echo.Context, userId string, params jcalendarsrv.GetUsersUserIdEventsParams) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	es, err := s.application.EventExtractor.GetEventsInInterval(ctx, uint(iid), params.From, params.Till)
	if err != nil {
		return err
	}

	outputEvents := make([]jcalendarsrv.OutputEvent, len(es), len(es))
	for idx, e := range es {
		participants := make([]jcalendarsrv.OutputUser, len(e.Users), len(e.Users))
		for uidx, p := range e.Users {
			participants[uidx] = jcalendarsrv.OutputUser{
				ID:             pcaster(int(p.ID)),
				CreateAt:       pcaster(p.CreatedAt.String()),
				UpdateAt:       pcaster(p.UpdatedAt.String()),
				FirstName:      &p.FirstName,
				LastName:       &p.LastName,
				Email:          &p.Email,
				TimeZoneOffset: &p.TimeZoneOffset,
			}
		}

		outputEvents[idx] = jcalendarsrv.OutputEvent{
			ID:       pcaster(int(e.ID)),
			CreateAt: pcaster(e.CreatedAt.String()),
			UpdateAt: pcaster(e.UpdatedAt.String()),
			From:     pcaster(e.From.String()),
			Till:     pcaster(e.Till.String()),
			Creator: &jcalendarsrv.OutputUser{
				ID:             pcaster(int(e.User.ID)),
				CreateAt:       pcaster(e.User.CreatedAt.String()),
				UpdateAt:       pcaster(e.User.UpdatedAt.String()),
				FirstName:      &e.User.FirstName,
				LastName:       &e.User.LastName,
				Email:          &e.User.Email,
				TimeZoneOffset: &e.User.TimeZoneOffset,
			},
			Participants: &participants,
			Details:      &e.Details,
			ScheduleRule: &e.ScheduleRule,
			IsPrivate:    &e.IsPrivate,
			IsRepeat:     &e.IsRepeat,
		}
	}

	return c.JSON(http.StatusOK, jcalendarsrv.EventsResponse{Data: &outputEvents})
}
