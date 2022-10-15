package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"jcalendar/internal/pkg"
	"jcalendar/internal/service/app"
	eevent "jcalendar/internal/service/entity/event"
	einvite "jcalendar/internal/service/entity/invite"
	euser "jcalendar/internal/service/entity/user"
	"jcalendar/internal/service/gateways/openapi/jcalendar"
	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/invites"
	"jcalendar/internal/service/repository/users"
	jcalendarsrv "jcalendar/pkg/openapi/jcalendar"
)

func TestPutInvitesId(t *testing.T) {
	var (
		ctx    = context.Background()
		path   = "/invtes/{id}"
		method = "PUT"
	)

	table := []*struct {
		testSubTittle       string
		inviteID            string
		inviteUpdateRequest jcalendarsrv.InviteUpdateRequest
		expectedStatus      int
	}{
		{
			testSubTittle: "invalid inviteID",
			inviteID:      "0",
			inviteUpdateRequest: jcalendarsrv.InviteUpdateRequest{
				Data: jcalendarsrv.InviteUpdate{
					IsAccepted: true,
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			testSubTittle: "successfully updating exist invite",
			inviteID:      "1",
			inviteUpdateRequest: jcalendarsrv.InviteUpdateRequest{
				Data: jcalendarsrv.InviteUpdate{
					IsAccepted: true,
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	db, err := pkg.NewDB(ctx)
	require.NoError(t, err)

	var (
		erepo = events.NewRepository(db)
		irepo = invites.NewRepository(db)
		urepo = users.NewRepository(db)
	)

	for _, row := range table {
		t.Run(row.testSubTittle, func(t *testing.T) {
			defer func() {
				require.NoError(t, clearTables(ctx, db))
			}()

			var application *app.Application
			application, err = app.NewApplication(ctx, db)
			require.NoError(t, err)

			participant := euser.NewUser(ctx, userFirstName, userLastName, userEmail, userhp, userTimeZoneOffset)
			err = urepo.CreateUser(ctx, participant)
			require.NoError(t, err)

			require.NoError(t,
				erepo.CreateEvent(ctx,
					eevent.NewEvent(
						ctx,
						eventFromTimestamp,
						eventTillTimestamp,
						1,
						[]uint{},
						eventDetails,
						nil,
						nil,
						nil,
						nil,
						false,
						false,
					),
				),
			)

			inv := einvite.NewInvite(ctx, inviteUserID, inviteEventID, false)

			_, err = irepo.CreateInvite(ctx, inv)
			require.NoError(t, err)

			var rawData []byte
			rawData, err = json.Marshal(row.inviteUpdateRequest)
			require.NoError(t, err)

			req := httptest.NewRequest(method, path, bytes.NewBuffer(rawData))
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(row.inviteID)
			c.Set("userID", uint(1))

			err = (&jcalendarsrv.ServerInterfaceWrapper{Handler: jcalendar.NewServer(application)}).PutInvitesId(c)
			if err != nil {
				actualErr := &echo.HTTPError{}
				require.True(t, errors.As(err, &actualErr))
				require.Equal(t, row.expectedStatus, actualErr.Code)
			} else {
				actual := jcalendarsrv.UpdatedInvite{}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual))

				var actualInvite *einvite.Invite
				actualInvite, err = irepo.GetInviteByID(ctx, uint(actual.ID))
				require.NoError(t, err)

				require.True(t, reflect.DeepEqual(
					einvite.Invite{
						ID:         actualInvite.ID,
						CreatedAt:  actualInvite.CreatedAt,
						UpdatedAt:  actualInvite.UpdatedAt,
						UserID:     1,
						EventID:    1,
						IsAccepted: true,
					},
					*actualInvite,
				))
			}
		})
	}
}
