package handler

import (
	"brevis/backend"
)

type (
	Handler struct {
		Backend backend.Backend
	}
)

// ErrorResponse holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}

// Response holds a response message
type Response struct {
	Message string `json:"message"`
}
