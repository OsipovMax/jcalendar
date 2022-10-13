package app

import (
	cmdevent "jcalendar/internal/service/usecase/commands/event"
	cmdinvite "jcalendar/internal/service/usecase/commands/invite"
	cmduser "jcalendar/internal/service/usecase/commands/user"
)

type Commands struct {
	CreateEvent cmdevent.CreateEventCommandHandler

	UpdateInvite cmdinvite.UpdateInviteCommandHandler

	CreateUser cmduser.CreateUserCommandHandler
}
