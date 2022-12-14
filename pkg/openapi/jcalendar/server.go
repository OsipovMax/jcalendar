// Package jcalendar provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package jcalendar

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Events information
	// (POST /events)
	PostEvents(ctx echo.Context) error
	// Event information
	// (GET /events/{id})
	GetEventsId(ctx echo.Context, id string) error
	// Updating status of invite
	// (PUT /invites/{id})
	PutInvitesId(ctx echo.Context, id string) error
	// Events information
	// (POST /login)
	PostLogin(ctx echo.Context) error
	// Adds user information
	// (POST /users)
	PostUsers(ctx echo.Context) error
	// Returns events information for user with user_id
	// (GET /users/{id}/events)
	GetUsersIdEvents(ctx echo.Context, id string, params GetUsersIdEventsParams) error
	// Returns closets free window for meeting
	// (GET /windows)
	GetWindows(ctx echo.Context, params GetWindowsParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostEvents converts echo context to params.
func (w *ServerInterfaceWrapper) PostEvents(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostEvents(ctx)
	return err
}

// GetEventsId converts echo context to params.
func (w *ServerInterfaceWrapper) GetEventsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEventsId(ctx, id)
	return err
}

// PutInvitesId converts echo context to params.
func (w *ServerInterfaceWrapper) PutInvitesId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutInvitesId(ctx, id)
	return err
}

// PostLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostLogin(ctx)
	return err
}

// PostUsers converts echo context to params.
func (w *ServerInterfaceWrapper) PostUsers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostUsers(ctx)
	return err
}

// GetUsersIdEvents converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersIdEvents(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUsersIdEventsParams
	// ------------- Required query parameter "from" -------------

	err = runtime.BindQueryParameter("form", true, true, "from", ctx.QueryParams(), &params.From)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter from: %s", err))
	}

	// ------------- Required query parameter "till" -------------

	err = runtime.BindQueryParameter("form", true, true, "till", ctx.QueryParams(), &params.Till)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter till: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUsersIdEvents(ctx, id, params)
	return err
}

// GetWindows converts echo context to params.
func (w *ServerInterfaceWrapper) GetWindows(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetWindowsParams
	// ------------- Required query parameter "user_ids[]" -------------

	err = runtime.BindQueryParameter("form", true, true, "user_ids[]", ctx.QueryParams(), &params.UserIds)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter user_ids[]: %s", err))
	}

	// ------------- Required query parameter "win_size" -------------

	err = runtime.BindQueryParameter("form", true, true, "win_size", ctx.QueryParams(), &params.WinSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter win_size: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWindows(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/events", wrapper.PostEvents)
	router.GET(baseURL+"/events/:id", wrapper.GetEventsId)
	router.PUT(baseURL+"/invites/:id", wrapper.PutInvitesId)
	router.POST(baseURL+"/login", wrapper.PostLogin)
	router.POST(baseURL+"/users", wrapper.PostUsers)
	router.GET(baseURL+"/users/:id/events", wrapper.GetUsersIdEvents)
	router.GET(baseURL+"/windows", wrapper.GetWindows)

}
