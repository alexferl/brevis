package db

import (
	"github.com/admiralobvious/brevis/internal/model"
)

// Database a common interface for all databases
type Database interface {
	Init() error
	Get(*model.UrlMapping) (*model.UrlMapping, error)
	GetStats(*model.UrlMapping) (*model.UrlMapping, error)
	Set(*model.UrlMapping) error
	Update(string, string, string) error
}
