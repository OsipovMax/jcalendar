package user

import (
	"context"
	"time"
)

type User struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	TimeZoneOffset int       `json:"time_zone_offset"`
	HashedPassword string    `json:"-"`
}

func NewUser(
	_ context.Context,
	firstName, lastName, email string,
	timeZoneOffset int,
) *User {
	return &User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		TimeZoneOffset: timeZoneOffset,
	}
}
