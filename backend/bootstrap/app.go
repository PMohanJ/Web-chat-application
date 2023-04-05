package bootstrap

import "github.com/pmohanj/web-chat-app/mongo"

type Application struct {
	Env   *Env
	Mongo mongo.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = DBinstance(app.Env.MongoDBURL)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseDBinstance(app.Mongo)
}
