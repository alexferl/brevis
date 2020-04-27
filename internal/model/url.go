package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/marksalpeter/token"

	"github.com/admiralobvious/brevis/internal/util"
)

type UrlMapping struct {
	Id             bson.ObjectId `json:"-" bson:"_id"`
	CreatedAt      *time.Time    `json:"created_at" bson:"created_at"`
	LastAccessedAt *time.Time    `json:"last_accessed_at" bson:"last_accessed_at"`
	ShortUrl       string        `json:"short_url" bson:"short_url"`
	UniqueViews    uint64        `json:"unique_views" bson:"unique_views"`
	Url            string        `json:"url" bson:"url"`
	UrlHash        [32]byte      `json:"-" bson:"url_hash"`
	Views          uint64        `json:"views" bson:"views"`
	Meta           `bson:"-"`
}

func NewUrlMapping(url string) *UrlMapping {
	now := time.Now().UTC()
	return &UrlMapping{
		Id:        bson.NewObjectId(),
		CreatedAt: &now,
		ShortUrl:  token.New().Encode(),
		Url:       url,
		UrlHash:   util.GetSha256Hash(url),
	}
}

type Referer struct {
	Id           bson.ObjectId `json:"-" bson:"_id"`
	Address      string        `json:"address" bson:"address"`
	AddressHash  [32]byte      `json:"-" bson:"address_hash"`
	FirstVisitAt *time.Time    `json:"first_visit_at" bson:"first_visit_at"`
	LastVisitAt  *time.Time    `json:"last_visit_at" bson:"last_visit_at"`
	ShortUrl     string        `json:"-" bson:"short_url"`
	Visits       uint64        `json:"visits" bson:"visits"`
}

type Meta struct {
	Id            bson.ObjectId `json:"-" bson:"_id"`
	LastUpdatedAt *time.Time    `json:"last_updated_at" bson:"last_updated_at"`
	Referrers     []Referer     `json:"referrers" bson:"referrers"`
	ShortUrl      string        `json:"-" bson:"short_url"`
	Visitors      [][32]byte    `json:"-" bson:"visitors"`
}
