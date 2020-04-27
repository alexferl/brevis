package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Root(c echo.Context) error {
	m := fmt.Sprint("Brevis: URL shortener API")
	return c.JSON(http.StatusOK, Response{Message: m})
}
