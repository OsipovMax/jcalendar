package jcalendar

import (
	"net/http"

	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"

	"github.com/labstack/echo/v4"
)

func (s *Server) GetUsersId(c echo.Context, id string) error {
	ctx := c.Request().Context()

	u, err := s.manager.GetUserByID(ctx, id)
	if err != nil {
		return err // TODO define problems + logs
	}

	return c.JSON(
		http.StatusOK,
		jcalendarsrv.UserReponse{
			Data: &jcalendarsrv.OutputUser{
				ID:             pcaster(int(u.ID)),
				CreateAt:       pcaster(u.CreatedAt.String()),
				UpdateAt:       pcaster(u.UpdatedAt.String()),
				FirstName:      &u.FirstName,
				LastName:       &u.LastName,
				Email:          &u.Email,
				TimeZoneOffset: &u.TimeZoneOffset,
			},
		},
	)
}
