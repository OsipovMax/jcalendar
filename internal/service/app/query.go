package app

import (
	qrevent "jcalendar/internal/service/usecase/queries/event"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"
)

type Queries struct {
	GetEvent qrevent.GetEventQueryHandler

	GetInvite qrinvite.GetInviteQueryHandler
}
