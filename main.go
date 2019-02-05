package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/srimaln91/geofence/app"
	"github.com/srimaln91/geofence/config"
)

func main() {

	//Read and parse config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	appConfig := &config.AppConfig{}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&appConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Println("config file has been loaded successfully.")

	//Initialize the app
	app.Init(appConfig)
}
