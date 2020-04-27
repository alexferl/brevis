package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/admiralobvious/brevis/internal/db"
)

type (
	// Handler represents the structure of our resource
	Handler struct {
		Database db.Database
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

// Register routes with echo
func Register(e *echo.Echo) {
	database := viper.Get("database").(db.Database)
	h := &Handler{Database: database}
	e.GET("/", h.Root)
	e.GET("/:id", h.Redirect)
	e.GET("/:id/stats", h.Stats)
	e.POST("/shorten", h.Shorten)
	e.POST("/unshorten", h.Unshorten)
}
