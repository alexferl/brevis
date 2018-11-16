package main

import (
	"strings"

	"github.com/admiralobvious/brevis/backend"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BackendFactory() backend.Backend {
	backendType := viper.GetString("backend-type")
	mongoDBTimeout := viper.GetDuration("backend-mongodb-timeout")
	mongoDBUri := viper.GetString("backend-mongodb-uri")
	logrus.Debugf("Using '%s' backend", backendType)

	switch strings.ToLower(backendType) {
	case "mongodb":
		return backend.NewMongoDBBackend(mongoDBUri, mongoDBTimeout)
	default:
		logrus.Warningf("Unknown backend type '%s'. Falling back to 'mongodb'", backendType)
		return backend.NewMongoDBBackend(mongoDBUri, mongoDBTimeout)
	}
}
