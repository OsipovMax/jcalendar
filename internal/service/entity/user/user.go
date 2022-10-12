package user

import (
	"context"
	"time"
)

type User struct {
	ID             uint      `json:"ID"`
	CreatedAt      time.Time `json:"CreatedAt"`
	UpdatedAt      time.Time `json:"UpdatedAt"`
	FirstName      string    `json:"FirstName"`
	LastName       string    `json:"LastName"`
	Email          string    `json:"Email"`
	TimeZoneOffset int       `json:"TimeZoneOffset"`
	HashedPassword string    `json:"-"`
}

func NewUser(
	_ context.Context,
	firstName, lastName, email, hashedPassword string,
	timeZoneOffset int,
) *User {
	return &User{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		HashedPassword: hashedPassword,
		TimeZoneOffset: timeZoneOffset,
	}
}
