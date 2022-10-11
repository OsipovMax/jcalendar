// Package jcalendar provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package jcalendar

// CreatedEvent defines model for CreatedEvent.
type CreatedEvent struct {
	ID *int `json:"ID,omitempty"`
}

// CreatedUser defines model for CreatedUser.
type CreatedUser struct {
	ID *int `json:"ID,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Status *int    `json:"Status,omitempty"`
	Title  *string `json:"Title,omitempty"`
}

// EventRequest defines model for EventRequest.
type EventRequest struct {
	Data *InputEvent `json:"data,omitempty"`
}

// EventResponse defines model for EventResponse.
type EventResponse struct {
	Data *OutputEvent `json:"data,omitempty"`
}

// EventsResponse defines model for EventsResponse.
type EventsResponse struct {
	Data *[]OutputEvent `json:"data,omitempty"`
}

// FreeWindow defines model for FreeWindow.
type FreeWindow struct {
	From *string `json:"From,omitempty"`
	Till *string `json:"Till,omitempty"`
}

// FreeWindowResponse defines model for FreeWindowResponse.
type FreeWindowResponse struct {
	Data *FreeWindow `json:"data,omitempty"`
}

// InputEvent defines model for InputEvent.
type InputEvent struct {
	CreatorID    *int    `json:"CreatorID,omitempty"`
	Details      *string `json:"Details,omitempty"`
	From         *string `json:"From,omitempty"`
	IsPrivate    *bool   `json:"IsPrivate,omitempty"`
	IsRepeat     *bool   `json:"IsRepeat,omitempty"`
	Participants *[]int  `json:"Participants,omitempty"`
	Till         *string `json:"Till,omitempty"`
}

// InputUser defines model for InputUser.
type InputUser struct {
	Email          *string `json:"Email,omitempty"`
	FirstName      *string `json:"FirstName,omitempty"`
	LastName       *string `json:"LastName,omitempty"`
	TimeZoneOffset *int    `json:"TimeZoneOffset,omitempty"`
}

// InviteUpdateRequest defines model for InviteUpdateRequest.
type InviteUpdateRequest struct {
	IsAccepted *bool `json:"IsAccepted,omitempty"`
}

// OutputEvent defines model for OutputEvent.
type OutputEvent struct {
	CreateAt  *string `json:"CreateAt,omitempty"`
	CreatorID *string `json:"CreatorID,omitempty"`
	Details   *string `json:"Details,omitempty"`
	From      *string `json:"From,omitempty"`
	ID        *int    `json:"ID,omitempty"`
	IsPrivate *bool   `json:"IsPrivate,omitempty"`
	Till      *string `json:"Till,omitempty"`
	UpdateAt  *string `json:"UpdateAt,omitempty"`
}

// OutputUser defines model for OutputUser.
type OutputUser struct {
	CreateAt       *string `json:"CreateAt,omitempty"`
	Email          *string `json:"Email,omitempty"`
	FirstName      *string `json:"FirstName,omitempty"`
	ID             *int    `json:"ID,omitempty"`
	LastName       *string `json:"LastName,omitempty"`
	TimeZoneOffset *int    `json:"TimeZoneOffset,omitempty"`
	UpdateAt       *string `json:"UpdateAt,omitempty"`
}

// UserReponse defines model for UserReponse.
type UserReponse struct {
	Data *OutputUser `json:"data,omitempty"`
}

// UserRequest defines model for UserRequest.
type UserRequest struct {
	Data *InputUser `json:"data,omitempty"`
}

// PostEventsJSONBody defines parameters for PostEvents.
type PostEventsJSONBody = EventRequest

// GetEventsUserIdParams defines parameters for GetEventsUserId.
type GetEventsUserIdParams struct {
	// begin of interval
	From string `form:"from" json:"from"`

	// end of interval
	Till string `form:"till" json:"till"`
}

// PutInvitesIdJSONBody defines parameters for PutInvitesId.
type PutInvitesIdJSONBody = InviteUpdateRequest

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody = UserRequest

// GetWindowsParams defines parameters for GetWindows.
type GetWindowsParams struct {
	// users identificators
	UserIds []string `form:"user_ids[]" json:"user_ids[]"`
}

// PostEventsJSONRequestBody defines body for PostEvents for application/json ContentType.
type PostEventsJSONRequestBody = PostEventsJSONBody

// PutInvitesIdJSONRequestBody defines body for PutInvitesId for application/json ContentType.
type PutInvitesIdJSONRequestBody = PutInvitesIdJSONBody

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody = PostUsersJSONBody
