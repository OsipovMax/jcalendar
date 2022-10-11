package jcalendar

import (
	"github.com/labstack/echo/v4"
)

func (s *Server) GetUsersEmail(c echo.Context, email string) error {
	//ctx := c.Request().Context()
	//
	//u, err := s.manager.GetUserByEmail(ctx, email)
	//if err != nil {
	//	return err // TODO define problems
	//}
	//
	//return c.JSON(
	//	http.StatusOK,
	//	jcalendarsrv.UserReponse{
	//		Data: &jcalendarsrv.OutputUser{
	//			ID:             pcaster(int(u.ID)),
	//			CreateAt:       pcaster(u.CreatedAt.String()),
	//			UpdateAt:       pcaster(u.UpdatedAt.String()),
	//			FirstName:      &u.FirstName,
	//			LastName:       &u.LastName,
	//			Email:          &u.Email,
	//			TimeZoneOffset: &u.TimeZoneOffset,
	//		},
	//	},
	//)
	return nil
}
