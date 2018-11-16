package handler

import (
	"github.com/admiralobvious/brevis/backend"
)

type (
	Handler struct {
		Backend backend.Backend
	}
)

// Error holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}
