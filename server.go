package main

import (
	"github.com/admiralobvious/brevis/handler"

	"github.com/admiralobvious/brevis/backend"
	"github.com/bakatz/echo-logrusmiddleware"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	cnf := NewConfig()
	cnf.BindFlags()

	InitLogging()
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.POST, echo.OPTIONS},
	}))

	if !viper.GetBool("log-requests-disabled") {
		e.Logger = logrusmiddleware.Logger{Logger: logrus.StandardLogger()}
		e.Use(logrusmiddleware.Hook())
	}

	b := viper.Get("backend").(backend.Backend)
	h := &handler.Handler{Backend: b}

	// Routes
	e.GET("/", h.Root)
	e.GET("/:id", h.Redirect)
	e.GET("/:id/stats", h.Stats)
	e.POST("/shorten", h.Shorten)
	e.POST("/unshorten", h.Unshorten)

	// Start server
	addr := viper.GetString("address") + ":" + viper.GetString("port")
	e.Logger.Fatal(e.Start(addr))
}
