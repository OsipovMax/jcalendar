package app

import (
	"context"

	"gorm.io/gorm"

	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/invites"
	"jcalendar/internal/service/repository/users"
	cmdevent "jcalendar/internal/service/usecase/commands/event"
	cmdinvite "jcalendar/internal/service/usecase/commands/invite"
	cmduser "jcalendar/internal/service/usecase/commands/user"
	mevent "jcalendar/internal/service/usecase/managers/event"
	mrule "jcalendar/internal/service/usecase/managers/rule"
	qrevent "jcalendar/internal/service/usecase/queries/event"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"
	qruser "jcalendar/internal/service/usecase/queries/user"
)

type Application struct {
	Commands     Commands
	Queries      Queries
	EventManager *mevent.EventManager
}

func NewApplication(ctx context.Context, db *gorm.DB) (*Application, error) {
	var (
		erepo = events.NewRepository(db)
		irepo = invites.NewRepository(db)
		urepo = users.NewRepository(db)
	)

	var (
		eventManager = mevent.NewEventManager(ctx, qrevent.NewGetEventsInIntervalQueryHandler(erepo))
		ruleManager  = mrule.NewRuleManager(ctx)
	)

	app := &Application{
		Commands: Commands{
			CreateEvent: cmdevent.NewCreateEventCommandHandler(erepo, ruleManager),

			UpdateInvite: cmdinvite.NewUpdateInviteCommandHandler(irepo),

			CreateUser: cmduser.NewCreateUserCommandHandler(urepo),
		},

		Queries: Queries{
			GetUserByEmail: qruser.NewGetUserByEmailQueryHandler(urepo),

			GetEvent: qrevent.NewGetEventQueryHandler(erepo),

			GetInvite: qrinvite.NewGetInviteQueryHandler(irepo),
		},

		EventManager: eventManager,
	}

	return app, nil
}
