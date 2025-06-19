package config

import "github.com/spf13/viper"

type MongoConfig struct {
	URI        string
	DBName     string
	Collection string
}

func newMongoConfig() *MongoConfig {
	return &MongoConfig{
		URI:        viper.GetString("MONGO_URI"),
		DBName:     viper.GetString("MONGO_DB_NAME"),
		Collection: viper.GetString("MONGO_EVENTS_COLLECTION"),
	}
}
