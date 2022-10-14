package app

import (
	"context"
	"os"

	"jcalendar/internal/service/repository/events"
	"jcalendar/internal/service/repository/invites"
	"jcalendar/internal/service/repository/users"
	cmdevent "jcalendar/internal/service/usecase/commands/event"
	cmdinvite "jcalendar/internal/service/usecase/commands/invite"
	cmduser "jcalendar/internal/service/usecase/commands/user"
	"jcalendar/internal/service/usecase/manager"
	qrevent "jcalendar/internal/service/usecase/queries/event"
	qrinvite "jcalendar/internal/service/usecase/queries/invite"
	qruser "jcalendar/internal/service/usecase/queries/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	postgresDSNKey = "POSTGRES"
)

type Application struct {
	Commands     Commands
	Queries      Queries
	EventManager *manager.EventManager
}

func NewApplication(ctx context.Context, db *gorm.DB) (*Application, error) {
	db, err := NewDB(ctx)
	if err != nil {
		return nil, err
	}

	var (
		erepo = events.NewRepository(db)
		irepo = invites.NewRepository(db)
		urepo = users.NewRepository(db)
	)

	m := manager.NewEventManager(ctx, qrevent.NewGetEventsInIntervalQueryHandler(erepo))

	app := &Application{
		Commands: Commands{
			CreateEvent: cmdevent.NewCreateEventCommandHandler(erepo, m),

			UpdateInvite: cmdinvite.NewUpdateInviteCommandHandler(irepo),

			CreateUser: cmduser.NewCreateUserCommandHandler(urepo),
		},

		Queries: Queries{
			GetUserByEmail: qruser.NewGetUserByEmailQueryHandler(urepo),

			GetEvent: qrevent.NewGetEventQueryHandler(erepo),

			GetInvite: qrinvite.NewGetInviteQueryHandler(irepo),
		},

		EventManager: m,
	}

	return app, nil
}

func NewDB(_ context.Context) (*gorm.DB, error) {
	var (
		dsn = os.Getenv(postgresDSNKey)
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
