package app

type Application struct {
	Commands Commands
	Queries  Queries
}

func NewApplication() *Application {
	app := &Application{
		Commands: Commands{},
		Queries:  Queries{},
	}

	return app
}
