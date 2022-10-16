package jcalendar

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	mevent "jcalendar/internal/service/usecase/managers/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

// nolint:revive,stylecheck
func (s *Server) GetUsersIdEvents(c echo.Context, id string, params jcalendarsrv.GetUsersIdEventsParams) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t convert string userID param for getting events in interval: %v", err)
		return echo.ErrBadRequest
	}

	es, err := s.application.EventManager.GetEventsInInterval(ctx, uint(iid), params.From, params.Till)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t get events in interval from event : %v", err)
		if errors.Is(err, mevent.ErrInvalidExecutingQuery) {
			if es == nil {
				return echo.ErrInternalServerError
			}
			return echo.ErrNotFound
		}
		return echo.ErrBadRequest
	}

	outputEvents := make([]jcalendarsrv.OutputEvent, len(es))
	for idx, e := range es {
		participants := make([]jcalendarsrv.OutputUser, len(e.Users))
		for uidx, p := range e.Users {
			participants[uidx] = jcalendarsrv.OutputUser{
				ID:             int(p.ID),
				CreatedAt:      p.CreatedAt.String(),
				UpdatedAt:      p.UpdatedAt.String(),
				FirstName:      p.FirstName,
				LastName:       p.LastName,
				Email:          p.Email,
				TimeZoneOffset: p.TimeZoneOffset,
			}
		}

		outputEvents[idx] = jcalendarsrv.OutputEvent{
			ID:        int(e.ID),
			CreatedAt: e.CreatedAt.String(),
			UpdatedAt: e.UpdatedAt.String(),
			From:      e.From.String(),
			Till:      e.Till.String(),
			Creator: jcalendarsrv.OutputUser{
				ID:             int(e.User.ID),
				CreatedAt:      e.User.CreatedAt.String(),
				UpdatedAt:      e.User.UpdatedAt.String(),
				FirstName:      e.User.FirstName,
				LastName:       e.User.LastName,
				Email:          e.User.Email,
				TimeZoneOffset: e.User.TimeZoneOffset,
			},
			Participants: &participants,
			Details:      e.Details,
			ScheduleRule: e.ScheduleRule,
			IsPrivate:    e.IsPrivate,
			IsRepeat:     e.IsRepeat,
		}
	}

	return c.JSON(http.StatusOK, jcalendarsrv.EventsResponse{Data: outputEvents})
}
