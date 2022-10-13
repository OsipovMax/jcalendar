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
	"jcalendar/internal/service/usecase/extractor"
	"jcalendar/internal/service/usecase/finder"
	"jcalendar/internal/service/usecase/parser"
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
	Commands       Commands
	Queries        Queries
	EventExtractor *extractor.EventExtractor
	Finder         *finder.Finder
}

func NewApplication(ctx context.Context) (*Application, error) {
	db, err := NewDB(ctx)
	if err != nil {
		return nil, err
	}

	var (
		erepo = events.NewRepository(db)
		irepo = invites.NewRepository(db)
		urepo = users.NewRepository(db)
	)

	ext := extractor.NewEventExtractor(ctx, qrevent.NewGetEventsInIntervalQueryHandler(erepo))

	app := &Application{
		Commands: Commands{
			CreateEvent: cmdevent.NewCreateEventCommandHandler(erepo, parser.NewScheduleParser(ctx)),

			CreateInvite: cmdinvite.NewCreateInviteCommandHandler(irepo),
			UpdateInvite: cmdinvite.NewUpdateInviteCommandHandler(irepo),

			CreateUser: cmduser.NewCreateUserCommandHandler(urepo),
		},

		Queries: Queries{
			GetUserByEmail: qruser.NewGetUserByEmailQueryHandler(urepo),

			GetEvent: qrevent.NewGetEventQueryHandler(erepo),

			GetInvite: qrinvite.NewGetInviteQueryHandler(irepo),
		},

		EventExtractor: ext,

		Finder: finder.NewFinder(ctx, ext),
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
