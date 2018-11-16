package backend

import (
	"errors"
	"fmt"
	"time"

	"github.com/admiralobvious/brevis/model"
	"github.com/admiralobvious/brevis/util"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
)

const (
	dbName   = "brevis"
	urlsColl = "urls"
)

type MongoDBBackend struct {
	Session *mgo.Session
	Timeout time.Duration
	Uri     string
}

func NewMongoDBBackend(uri string, timeout time.Duration) Backend {
	return Backend(&MongoDBBackend{
		Timeout: timeout,
		Uri:     uri,
	})
}

func (mb *MongoDBBackend) Init() (err error) {
	b := &backoff.Backoff{
		Jitter: true,
	}

	for {
		mb.Session, err = mgo.DialWithTimeout(mb.Uri, mb.Timeout)
		if err != nil {
			d := b.Duration()
			logrus.Errorf("%s, reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}

		b.Reset()

		coll := mb.Session.DB(dbName).C(urlsColl)
		if coll == nil {
			m := fmt.Sprint("Error creating collection")
			return errors.New(m)
		}

		urlIndex := mgo.Index{
			Key:      []string{"url_hash"},
			Unique:   true,
			DropDups: true,
		}
		uErr := coll.EnsureIndex(urlIndex)
		if uErr != nil {
			return uErr
		}

		shortUrlIndex := mgo.Index{
			Key:      []string{"short_url"},
			Unique:   true,
			DropDups: true,
		}
		sErr := coll.EnsureIndex(shortUrlIndex)
		if sErr != nil {
			return sErr
		}

		return nil
	}
}

func (mb *MongoDBBackend) Get(mapping *model.UrlMapping) (*model.UrlMapping, error) {
	session := mb.Session.Copy()
	defer session.Close()

	result := model.UrlMapping{}

	pipeline := bson.M{
		"$or": []interface{}{
			bson.M{"url_hash": util.GetSha256Hash(mapping.Url)},
			bson.M{"short_url": mapping.ShortUrl},
		},
	}

	err := session.DB(dbName).C(urlsColl).Find(pipeline).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		logrus.Errorf("Error searching: %s", err)
		return nil, err
	}

	return &result, nil
}

func (mb *MongoDBBackend) Set(mapping *model.UrlMapping) error {
	session := mb.Session.Copy()
	defer session.Close()

	err := session.DB(dbName).C(urlsColl).Insert(mapping)
	if err != nil {
		if mgo.IsDup(err) {
			res, err := mb.Get(mapping)
			if err != nil {
				logrus.Errorf("Error getting result: %s", err)
				return err
			}
			mapping.ShortUrl = res.ShortUrl
			return nil
		}
		logrus.Errorf("Error inserting: %s", err)
		return err
	}

	return nil
}

func (mb *MongoDBBackend) Update(mapping *model.UrlMapping) error {
	session := mb.Session.Copy()
	defer session.Close()

	var result bson.M
	change := mgo.Change{
		Update: bson.M{"$inc": bson.M{"views": 1}, "$set": bson.M{"last_accessed_at": time.Now().UTC()}},
	}
	_, err := session.DB(dbName).C(urlsColl).Find(bson.M{"short_url": mapping.ShortUrl}).Apply(change, &result)

	if err != nil {
		logrus.Errorf("Error updating: %s", err)
		return err
	}

	return nil
}
