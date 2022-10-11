package app

import (
	qrevent "jcalendar/internal/service/usecase/queries/event"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"
	qruser "jcalendar/internal/service/usecase/queries/user"
)

type Queries struct {
	GetUserByEmail qruser.GetUserByEmailQueryHandler

	GetEvent qrevent.GetEventQueryHandler

	GetInvite qrinvite.GetInviteQueryHandler
}
