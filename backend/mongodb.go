package backend

import (
	"errors"
	"fmt"
	"time"

	"github.com/admiralobvious/brevis/model"

	"github.com/Sirupsen/logrus"
	"github.com/jpillora/backoff"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDBBackend struct {
	Session *mgo.Session
	Uri     string
}

func NewMongoDBBackend(uri string) Backend {
	return Backend(&MongoDBBackend{
		Uri: uri,
	})
}

func (mb *MongoDBBackend) Init() (err error) {
	b := &backoff.Backoff{
		Jitter: true,
	}

	for {
		mb.Session, err = mgo.Dial(mb.Uri)
		mb.Session.SetSocketTimeout(time.Second * 3)
		mb.Session.SetSyncTimeout(time.Second * 3)

		if err != nil {
			d := b.Duration()
			logrus.Errorf("%s, reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}

		b.Reset()

		coll := mb.Session.DB("brevis").C("urls")
		if coll == nil {
			m := fmt.Sprint("Error creating collection")
			logrus.Error(m)
			return errors.New(m)
		}

		urlIndex := mgo.Index{
			Key:      []string{"url"},
			Unique:   true,
			DropDups: true,
		}
		coll.EnsureIndex(urlIndex)
		shortUrlIndex := mgo.Index{
			Key:      []string{"shortUrl"},
			Unique:   true,
			DropDups: true,
		}
		coll.EnsureIndex(shortUrlIndex)

		return nil
	}
}

func (mb *MongoDBBackend) Get(mapping *model.UrlMapping) (*model.UrlMapping, error) {
	result := model.UrlMapping{}
	session := mb.Session.Copy()
	defer session.Close()

	pipeline := bson.M{
		"$or": []interface{}{
			bson.M{"url": mapping.Url},
			bson.M{"shortUrl": mapping.ShortUrl},
		},
	}

	err := session.DB("brevis").C("urls").Find(pipeline).One(&result)

	if err != nil && err != mgo.ErrNotFound {
		logrus.Errorf("Error searching: %s", err)
		return nil, err
	}

	return &result, nil
}

func (mb *MongoDBBackend) Set(mapping *model.UrlMapping) error {
	session := mb.Session.Copy()
	defer session.Close()

	err := session.DB("brevis").C("urls").Insert(mapping)
	if err != nil {
		if mgo.IsDup(err) {
			res, err := mb.Get(mapping)
			if err != nil {
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
