package user

import (
	"context"
	"errors"
)

type CreateUserCommand struct {
	FirstName      string
	LastName       string
	Email          string
	TimeZoneOffset int
	HashedPassword string
}

func NewCreateUserCommand(
	_ context.Context,
	firstName, lastName, email, hashedPassword string,
	timeZoneOffset int,
) (*CreateUserCommand, error) {
	if firstName == "" {
		return nil, errors.New("missing firstName value")
	}

	if lastName == "" {
		return nil, errors.New("missing lastName value")
	}

	if email == "" {
		return nil, errors.New("missing email value")
	}

	if hashedPassword == "" {
		return nil, errors.New("missing password value")
	}

	if timeZoneOffset < 0 {
		return nil, errors.New("negative timeZoneOffset value")
	}

	return &CreateUserCommand{
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		TimeZoneOffset: timeZoneOffset,
		HashedPassword: hashedPassword,
	}, nil
}
