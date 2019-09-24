package backend

import (
	"errors"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"

	"brevis/model"
	"brevis/util"
)

const (
	dbName   = "brevis"
	metaColl = "meta"
	refsColl = "referrers"
	urlsColl = "urls"
)

type MongoDBBackend struct {
	Username string
	Password string
	Session *mgo.Session
	Timeout time.Duration
	Uri     string
}

func NewMongoDBBackend(uri string, timeout time.Duration, username, password string) Backend {
	return Backend(&MongoDBBackend{
		Username: username,
		Password: password,
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
		if mb.Username != "" && mb.Password != "" {
			creds := mgo.Credential{Username: mb.Username, Password: mb.Password}
			err = mb.Session.Login(&creds)
		}
		if err != nil {
			d := b.Duration()
			logrus.Errorf("%s, reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}

		b.Reset()

		err := mb.createIndexes()
		if err != nil {
			return err
		}

		return nil
	}
}

func (mb *MongoDBBackend) Get(mapping *model.UrlMapping) (*model.UrlMapping, error) {
	session := mb.Session.Copy()
	defer session.Close()

	result := model.UrlMapping{}
	err := session.DB(dbName).C(urlsColl).Find(bson.M{"short_url": mapping.ShortUrl}).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		logrus.Errorf("Error searching: %s", err)
		return nil, err
	}

	return &result, nil
}

func (mb *MongoDBBackend) GetStats(mapping *model.UrlMapping) (*model.UrlMapping, error) {
	session := mb.Session.Copy()
	defer session.Close()

	db := session.DB(dbName)
	query := bson.M{"short_url": mapping.ShortUrl}

	var result model.UrlMapping
	err := db.C(urlsColl).Find(query).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		logrus.Errorf("Error searching url mapping: %s", err)
		return nil, err
	}

	var meta model.Meta
	err = db.C(metaColl).Find(query).One(&meta)
	if err != nil && err != mgo.ErrNotFound {
		logrus.Errorf("Error searching meta: %s", err)
		return nil, err
	}

	var referrers []model.Referer
	err = db.C(refsColl).Find(query).All(&referrers)
	if err != nil && err != mgo.ErrNotFound {
		logrus.Errorf("Error searching referrers: %s", err)
		return nil, err
	}

	meta.Referrers = referrers
	result.Meta = meta

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
		logrus.Errorf("Error inserting url mapping: %s", err)
		return err
	}

	return nil
}

func (mb *MongoDBBackend) Update(shortUrl, referer, visitor string) error {
	session := mb.Session.Copy()
	defer session.Close()

	db := session.DB(dbName)
	query := bson.M{"short_url": shortUrl}

	var doc model.UrlMapping
	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{"views": 1},
			"$set": bson.M{"last_accessed_at": time.Now().UTC()}},
		ReturnNew: true,
	}
	_, err := db.C(urlsColl).Find(query).Apply(change, &doc)
	if err != nil {
		logrus.Errorf("Error updating url mapping: %s", err)
		return err
	}

	var ref model.Referer
	change = mgo.Change{
		Update: bson.M{
			"$inc": bson.M{"visits": 1},
			"$set": bson.M{"last_visit_at": time.Now().UTC()},
			"$setOnInsert": bson.M{
				"first_visit_at": time.Now().UTC(),
				"address":        referer,
				"address_hash":   util.GetSha256Hash(referer)}},
		Upsert: true,
	}
	_, err = db.C(refsColl).Find(bson.M{
		"short_url":    shortUrl,
		"address_hash": util.GetSha256Hash(referer)}).Apply(change, &ref)
	if err != nil {
		logrus.Errorf("Error updating referer: %s", err)
		return err
	}

	var meta model.Meta
	change = mgo.Change{
		Update: bson.M{
			"$addToSet": bson.M{"visitors": util.GetSha256Hash(visitor)},
			"$set":      bson.M{"last_updated_at": time.Now().UTC()}},
		Upsert:    true,
		ReturnNew: true,
	}
	_, err = db.C(metaColl).Find(query).Apply(change, &meta)
	if err != nil {
		logrus.Errorf("Error updating meta: %s", err)
		return err
	}

	change = mgo.Change{
		Update:    bson.M{"$set": bson.M{"unique_views": uint64(len(meta.Visitors))}},
		ReturnNew: true,
	}
	_, err = db.C(urlsColl).Find(query).Apply(change, &doc)
	if err != nil {
		logrus.Errorf("Error updating unique views: %s", err)
		return err
	}

	return nil
}

func (mb *MongoDBBackend) createIndexes() error {
	session := mb.Session.Copy()
	defer session.Close()

	db := mb.Session.DB(dbName)

	// urls
	urlColl := db.C(urlsColl)
	if urlColl == nil {
		m := fmt.Sprint("Error creating urls collection")
		return errors.New(m)
	}

	err := urlColl.EnsureIndex(mgo.Index{
		Key:      []string{"short_url"},
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		return err
	}

	// meta
	metaColl := db.C(metaColl)
	if metaColl == nil {
		m := fmt.Sprint("Error creating meta collection")
		return errors.New(m)
	}

	err = metaColl.EnsureIndex(mgo.Index{
		Key:      []string{"short_url"},
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		return err
	}

	// referrers
	refColl := db.C(refsColl)
	if refColl == nil {
		m := fmt.Sprint("Error creating referrers collection")
		return errors.New(m)
	}

	err = refColl.EnsureIndex(mgo.Index{
		Key:      []string{"short_url", "address_hash"},
		Unique:   true,
		DropDups: true,
	})
	if err != nil {
		return err
	}

	return nil
}
