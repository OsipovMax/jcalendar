package jcalendar

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	qrevent "jcalendar/internal/service/usecase/queries/event"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

// nolint:revive,stylecheck
func (s *Server) GetEventsId(c echo.Context, id string) error {
	ctx := c.Request().Context()

	iid, err := strconv.Atoi(id)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t convert string id param for getting events: %v", err)
		return echo.ErrBadRequest
	}

	query, err := qrevent.NewGetEventQuery(ctx, uint(iid))
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t create GetEventQuery: %v", err)
		return echo.ErrBadRequest
	}

	e, err := s.application.Queries.GetEvent.Handle(ctx, query)
	if err != nil {
		logrus.WithContext(ctx).Errorf("can`t execute GetEventQuery: %v", err)
		if e == nil {
			return echo.ErrInternalServerError
		}
		return echo.ErrNotFound
	}

	var details = e.Details
	if !isResourceOwner(ctx, e.CreatorID, c.Get(userIDClaim).(uint)) {
		details = busyEventDetail
	}

	participants := make([]jcalendarsrv.OutputUser, len(e.Users))
	for idx, participant := range e.Users {
		participants[idx] = jcalendarsrv.OutputUser{
			ID:             int(participant.ID),
			CreatedAt:      participant.CreatedAt.String(),
			UpdatedAt:      participant.UpdatedAt.String(),
			FirstName:      participant.FirstName,
			LastName:       participant.LastName,
			Email:          participant.Email,
			TimeZoneOffset: participant.TimeZoneOffset,
		}
	}

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.EventResponse{
			Data: jcalendarsrv.OutputEvent{
				ID:           int(e.ID),
				CreatedAt:    e.CreatedAt.String(),
				UpdatedAt:    e.UpdatedAt.String(),
				From:         e.From.String(),
				Till:         e.Till.String(),
				Details:      details,
				ScheduleRule: e.ScheduleRule,
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
				IsPrivate:    e.IsPrivate,
				IsRepeat:     e.IsRepeat,
			},
		},
	)
}
