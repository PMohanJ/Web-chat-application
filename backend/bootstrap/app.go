package bootstrap

import (
	"os"

	"github.com/pmohanj/web-chat-app/mongo"
)

type Application struct {
	Env   *Env
	Mongo mongo.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()

	// To facilitate for testing env in github actions
	MongoDBURL := os.Getenv("MongoDBURL")
	if MongoDBURL == "" {
		app.Mongo = DBinstance(app.Env.MongoDBURL)
	} else {
		app.Mongo = DBinstance(MongoDBURL)
	}

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseDBinstance(app.Mongo)
}
