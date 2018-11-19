package backend

import (
	"github.com/admiralobvious/brevis/model"
)

// Backend a common interface for all backends
type Backend interface {
	Init() error
	Get(*model.UrlMapping) (*model.UrlMapping, error)
	GetStats(*model.UrlMapping) (*model.UrlMapping, error)
	Set(*model.UrlMapping) error
	Update(string, string, string) error
}
