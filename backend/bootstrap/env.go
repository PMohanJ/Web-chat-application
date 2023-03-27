package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	Origin            string `mapstructure:"ORIGIN"`
	ContextTimeout    string `mapstructure:"CONTEXT_TIMEOUT"`
	MongoDBURL        string `mapstructure:"MONGODB_URL"`
	MongoDBURLTesting string `mapstructure:"MONGODB_URL_TESTING"`
	SecretKey         string `mapstructure:"SECRET_KEY"`
	Port              string `mapstructure:"PORT"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Can't find the .env file: ", err)
	}

	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalln("Environment can't be loaded: ", err)
	}

	return &env
}
