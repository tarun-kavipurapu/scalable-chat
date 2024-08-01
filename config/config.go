package config

import (
	"fmt"
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

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigType("env")

	// Set the config file name
	viper.SetConfigName("app")

	// Add the path to the config file
	viper.AddConfigPath(path)

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			return config, fmt.Errorf("config file not found in %s: %w", path, err)
		}
		// Config file was found but another error was produced
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Attempt to unmarshal the config into the struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Print the full path of the config file that was successfully loaded
	log.Printf("Config file loaded successfully: %s", viper.ConfigFileUsed())

	return config, nil
}
