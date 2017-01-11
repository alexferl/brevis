package main

import (
	"strings"

	"github.com/admiralobvious/brevis/backend"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

func BackendFactory() backend.Backend {
	backendType := viper.GetString("backend-type")
	mongoDBUri := viper.GetString("backend-mongodb-uri")
	log.Debugf("Using '%s' backend", backendType)

	switch strings.ToLower(backendType) {
	case "mongodb":
		return backend.NewMongoDBBackend(mongoDBUri)
	default:
		log.Warningf("Unknown loader type '%s'. Falling back to 'MongoDB'", backendType)
		return backend.NewMongoDBBackend(mongoDBUri)
	}
}
