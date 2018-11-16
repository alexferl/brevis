package handler

import (
	"fmt"
	"net/http"

	"github.com/admiralobvious/brevis/model"
	"github.com/admiralobvious/brevis/util"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func (h *Handler) Shorten(c echo.Context) error {
	um := &model.UrlMapping{}
	if err := c.Bind(um); err != nil {
		m := fmt.Sprintf("%s", err.(*echo.HTTPError).Message)
		return c.JSON(err.(*echo.HTTPError).Code, ErrorResponse{Message: m})
	}

	if um.Url == "" {
		m := fmt.Sprintf("Field 'url' is required")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: m})
	}

	if len(um.Url) > 2048 {
		m := fmt.Sprintf("url is too long, must be 2048 characters or less")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: m})
	}

	valid := util.IsValidUri(um.Url)
	if !valid {
		m := fmt.Sprintf("url '%s' is not valid", um.Url)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: m})
	}

	su := model.NewShortUrl(um.Url)
	err := h.Backend.Set(su)
	if err != nil {
		m := fmt.Sprintf("Error shortening url")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: m})
	}

	baseUrl := viper.GetString("base-url")

	return c.JSON(http.StatusCreated, map[string]string{"short_url": baseUrl + su.ShortUrl})
}

func (h *Handler) Unshorten(c echo.Context) error {
	um := &model.UrlMapping{}
	if err := c.Bind(um); err != nil {
		m := fmt.Sprintf("%s", err.(*echo.HTTPError).Message)
		return c.JSON(err.(*echo.HTTPError).Code, ErrorResponse{Message: m})
	}

	if um.ShortUrl == "" {
		m := fmt.Sprintf("Field 'short_url' required")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: m})
	}

	res, err := h.Backend.Get(um)
	if err != nil {
		m := fmt.Sprintf("Error getting short_url")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: m})
	}

	if res.ShortUrl == "" {
		m := fmt.Sprintf("short_url '%s' not found", um.ShortUrl)
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: m})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": res.Url})
}

func (h *Handler) Redirect(c echo.Context) error {
	id := c.Param("id")
	um := &model.UrlMapping{ShortUrl: id}

	res, err := h.Backend.Get(um)
	if err != nil {
		m := fmt.Sprintf("Error getting id '%s'", id)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: m})
	}

	if res.Url == "" {
		m := fmt.Sprintf("id '%s' not found", id)
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: m})
	}

	uErr := h.Backend.Update(um)
	if uErr != nil {
		m := fmt.Sprintf("Error updating id '%s'", id)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: m})
	}

	return c.Redirect(301, res.Url)
}

func (h *Handler) Stats(c echo.Context) error {
	id := c.Param("id")
	um := &model.UrlMapping{ShortUrl: id}

	res, err := h.Backend.Get(um)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	if res.ShortUrl == "" {
		m := fmt.Sprintf("id '%s' not found", id)
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: m})
	}

	baseUrl := viper.GetString("base-url")
	res.ShortUrl = baseUrl + res.ShortUrl

	return c.JSON(http.StatusOK, res)
}
