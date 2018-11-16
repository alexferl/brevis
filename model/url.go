package model

import (
	"time"

	"github.com/admiralobvious/brevis/util"
)

type UrlMapping struct {
	CreatedAt      *time.Time `json:"created_at" bson:"created_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at" bson:"last_accessed_at"`
	ShortUrl       string     `json:"short_url" bson:"short_url"`
	Url            string     `json:"url"`
	UrlHash        [32]byte   `json:"-" bson:"url_hash"`
	Views          uint64     `json:"views" bson:"views"`
}

func NewShortUrl(url string) *UrlMapping {
	now := time.Now().UTC()
	return &UrlMapping{
		CreatedAt:      &now,
		LastAccessedAt: nil,
		ShortUrl:       NewToken(),
		Url:            url,
		UrlHash:        util.GetSha256Hash(url),
	}
}
