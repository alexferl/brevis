package main

import (
	"fmt"
	"net/http"

	"github.com/admiralobvious/brevis/backend"
	"github.com/admiralobvious/brevis/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func root(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Brevis: URL shortener API"})
}

func shorten(c echo.Context) error {
	b := viper.Get("backend").(backend.Backend)
	um := &model.UrlMapping{}
	if err := c.Bind(um); err != nil {
		m := fmt.Sprint("Must provide a URL to shorten")
		return c.JSON(http.StatusBadRequest, Error{m})
	}

	valid := IsValidUri(um.Url)
	if !valid {
		m := fmt.Sprintf("URL '%s' is not valid", um.Url)
		return c.JSON(http.StatusBadRequest, Error{m})
	}

	su := model.NewShortUrl(um.Url)
	err := b.Set(su)
	if err != nil {
		m := fmt.Sprintf("Error shortening URL '%s': %s", um.Url, err)
		return c.JSON(http.StatusInternalServerError, Error{m})
	}

	baseUrl := viper.GetString("base-url")
	su.ShortUrl = baseUrl + su.ShortUrl
	return c.JSON(http.StatusOK, su)

}

func unshorten(c echo.Context) error {
	b := viper.Get("backend").(backend.Backend)
	um := &model.UrlMapping{}
	if err := c.Bind(um); err != nil {
		m := fmt.Sprint("Must provide an id to unshorten")
		return c.JSON(http.StatusBadRequest, Error{m})
	}

	res, err := b.Get(um)
	if err != nil {
		m := fmt.Sprintf("id '%s' not found", um.ShortUrl)
		return c.JSON(http.StatusNotFound, Error{m})
	}

	baseUrl := viper.GetString("base-url")
	res.ShortUrl = baseUrl + res.ShortUrl
	return c.JSON(http.StatusOK, res)
}

func redirect(c echo.Context) error {
	id := c.Param("id")
	b := viper.Get("backend").(backend.Backend)

	um := &model.UrlMapping{}
	um.ShortUrl = id

	res, err := b.Get(um)
	if err != nil {
		m := fmt.Sprintf("Error getting id '%s': %v", id, err)
		return c.JSON(http.StatusInternalServerError, Error{m})
	}

	if res.Url == "" {
		m := fmt.Sprintf("id '%s' not found", id)
		return c.JSON(http.StatusOK, Error{m})
	}

	baseUrl := viper.GetString("base-url")
	res.ShortUrl = baseUrl + res.ShortUrl
	return c.Redirect(301, res.Url)
}

func init() {
	cnf := NewConfig()
	cnf.BindFlags()

	InitLogging()
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.POST},
	}))

	// Routes
	e.GET("/", root)
	e.GET("/:id", redirect)
	e.POST("/shorten", shorten)
	e.POST("/unshorten", unshorten)

	addr := viper.GetString("address") + ":" + viper.GetString("port")
	// Start server
	e.Logger.Fatal(e.Start(addr))
}
