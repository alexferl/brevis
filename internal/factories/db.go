package factories

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/admiralobvious/brevis/internal/db"
)

func DatabaseFactory() db.Database {
	databaseType := viper.GetString("database-type")
	mongoDBUsername := viper.GetString("database-mongodb-username")
	mongoDBPassword := viper.GetString("database-mongodb-password")
	mongoDBTimeout := viper.GetDuration("database-mongodb-timeout")
	mongoDBUri := viper.GetString("database-mongodb-uri")
	logrus.Debugf("Using '%s' database", databaseType)

	switch strings.ToLower(databaseType) {
	case "mongodb":
		return db.NewMongoDatabase(mongoDBUri, mongoDBTimeout, mongoDBUsername, mongoDBPassword)
	default:
		logrus.Warningf("Unknown database type '%s'. Falling back to 'mongodb'", databaseType)
		return db.NewMongoDatabase(mongoDBUri, mongoDBTimeout, mongoDBUsername, mongoDBPassword)
	}
}
