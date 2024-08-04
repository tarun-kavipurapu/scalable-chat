package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource             string `mapstructure:"DB_SOURCE"`
	DBDriver             string `mapstructure:"DB_DRIVER"`
	PORT                 string `mapstructure:"PORT"`
	AWS_REGION           string `mapstructure:"AWS_REGION"`
	AWS_ACCESS_TOKEN     string `mapstructure:"AWS_ACCESS_TOKEN"`
	AWS_SECRET_TOKEN_KEY string `mapstructure:"AWS_SECRET_TOKEN_KEY"`
	AWS_BUCKET_NAME      string `mapstructure:"AWS_BUCKET_NAME"`
}

var EnvVars Config

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigType("env")

	// Load app.env

	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	log.Println(path)
	err = viper.MergeInConfig()
	if err != nil {
		log.Println("Cannot read app config file:", err)
	}

	// Load keys.env
	viper.SetConfigName("keys")
	viper.AddConfigPath(path)
	err = viper.MergeInConfig()
	if err != nil {
		log.Println("Cannot read keys config file:", err)
	}

	viper.AutomaticEnv()

	EnvVars.PORT = viper.GetString("PORT")
	EnvVars.DBSource = viper.GetString("DB_SOURCE")
	EnvVars.DBDriver = viper.GetString("DB_DRIVER")
	EnvVars.AWS_REGION = viper.GetString("AWS_REGION")
	EnvVars.AWS_ACCESS_TOKEN = viper.GetString("AWS_ACCESS_TOKEN")
	EnvVars.AWS_SECRET_TOKEN_KEY = viper.GetString("AWS_SECRET_TOKEN_KEY")
	EnvVars.AWS_BUCKET_NAME = viper.GetString("AWS_BUCKET_NAME")

	return EnvVars, nil
}
