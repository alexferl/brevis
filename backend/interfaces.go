package backend

import (
	"github.com/admiralobvious/brevis/model"
)

// Backend a common interface for all backend
type Backend interface {
	Init() error
	Get(*model.UrlMapping) (*model.UrlMapping, error)
	Set(*model.UrlMapping) error
}
