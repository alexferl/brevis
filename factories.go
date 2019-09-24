package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"brevis/backend"
)

func BackendFactory() backend.Backend {
	backendType := viper.GetString("backend-type")
	mongoDBUsername := viper.GetString("backend-mongodb-username")
	mongoDBPassword := viper.GetString("backend-mongodb-password")
	mongoDBTimeout := viper.GetDuration("backend-mongodb-timeout")
	mongoDBUri := viper.GetString("backend-mongodb-uri")
	logrus.Debugf("Using '%s' backend", backendType)

	switch strings.ToLower(backendType) {
	case "mongodb":
		return backend.NewMongoDBBackend(mongoDBUri, mongoDBTimeout, mongoDBUsername, mongoDBPassword)
	default:
		logrus.Warningf("Unknown backend type '%s'. Falling back to 'mongodb'", backendType)
		return backend.NewMongoDBBackend(mongoDBUri, mongoDBTimeout, mongoDBUsername, mongoDBPassword)
	}
}
