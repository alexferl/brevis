package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func (h *Handler) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Brevis: URL shortener API"})
}
