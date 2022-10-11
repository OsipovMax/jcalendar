package app

import (
	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/invites"
	"jcalendar/internal/service/repository/users"
	cmdevent "jcalendar/internal/service/usecase/commands/event"
	cmdinvite "jcalendar/internal/service/usecase/commands/invite"
	cmduser "jcalendar/internal/service/usecase/commands/user"
	qrevent "jcalendar/internal/service/usecase/queries/event"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"
	qruser "jcalendar/internal/service/usecase/queries/user"

	"gorm.io/gorm"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

func NewApplication(db *gorm.DB) *Application {
	var (
		erepo = events.NewRepository(db)
		irepo = invites.NewRepository(db)
		urepo = users.NewRepository(db)
	)
	app := &Application{
		Commands: Commands{
			CreateEvent: cmdevent.NewCreateEventCommandHandler(erepo),

			CreateInvite: cmdinvite.NewCreateInviteCommandHandler(irepo),
			UpdateInvite: cmdinvite.NewUpdateInviteCommandHandler(irepo),

			CreateUser: cmduser.NewCreateUserCommandHandler(urepo),
		},

		Queries: Queries{
			GetUserByEmail: qruser.NewGetUserByEmailQueryHandler(urepo),

			GetEvent: qrevent.NewGetEventQueryHandler(erepo),

			GetInvite: qrinvite.NewGetInviteQueryHandler(irepo),
		},
	}

	return app
}
