package model

type UrlMapping struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shortUrl,omitempty"bson:"shortUrl"`
}

func NewShortUrl(url string) *UrlMapping {
	return &UrlMapping{
		Url:      url,
		ShortUrl: NewToken(),
	}
}
