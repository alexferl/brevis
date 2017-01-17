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
		m := fmt.Sprint("Must provide a URL to shorten")
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{Message: m})
	}

	valid := util.IsValidUri(um.Url)
	if !valid {
		m := fmt.Sprintf("URL '%s' is not valid", um.Url)
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{Message: m})
	}

	su := model.NewShortUrl(um.Url)
	err := h.Backend.Set(su)
	if err != nil {
		m := fmt.Sprintf("Error shortening URL '%s': %s", um.Url, err)
		return c.JSON(http.StatusInternalServerError, util.ErrorResponse{Message: m})
	}

	baseUrl := viper.GetString("base-url")
	su.ShortUrl = baseUrl + su.ShortUrl
	return c.JSON(http.StatusOK, su)

}

func (h *Handler) Unshorten(c echo.Context) error {
	um := &model.UrlMapping{}
	if err := c.Bind(um); err != nil {
		m := fmt.Sprint("Must provide an id to unshorten")
		return c.JSON(http.StatusBadRequest, util.ErrorResponse{Message: m})
	}

	res, err := h.Backend.Get(um)
	if err != nil {
		m := fmt.Sprintf("id '%s' not found", um.ShortUrl)
		return c.JSON(http.StatusNotFound, util.ErrorResponse{Message: m})
	}

	baseUrl := viper.GetString("base-url")
	res.ShortUrl = baseUrl + res.ShortUrl
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Redirect(c echo.Context) error {
	id := c.Param("id")

	um := &model.UrlMapping{}
	um.ShortUrl = id

	res, err := h.Backend.Get(um)
	if err != nil {
		m := fmt.Sprintf("Error getting id '%s': %v", id, err)
		return c.JSON(http.StatusInternalServerError, util.ErrorResponse{Message: m})
	}

	if res.Url == "" {
		m := fmt.Sprintf("id '%s' not found", id)
		return c.JSON(http.StatusNotFound, util.ErrorResponse{Message: m})
	}

	baseUrl := viper.GetString("base-url")
	res.ShortUrl = baseUrl + res.ShortUrl
	return c.Redirect(301, res.Url)
}
