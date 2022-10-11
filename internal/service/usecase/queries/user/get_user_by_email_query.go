package user

import "errors"

type GetUserByEmailQuery struct {
	UserEmail string
}

func NewUserByEmailQuery(userEmail string) (*GetUserByEmailQuery, error) {
	if userEmail == "" {
		return nil, errors.New("missing userEmail value")
	}

	return &GetUserByEmailQuery{
		UserEmail: userEmail,
	}, nil
}
