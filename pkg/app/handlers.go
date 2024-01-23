package app

import (
	"fmt"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *App) runHTTPServer(host string, port int) error {
	listenAddress := fmt.Sprintf("%s:%d", host, port)
	a.Printf("starting http listener at http://%s\n", listenAddress)

	return a.echo.Start(listenAddress)
}

func (a *App) registerHandlers() {
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{"Authorization", "Authorization2", "Origin", "X-Requested-With", "Content-Type", "Accept", "Platform", "Version"},
	}))

	// sentry middleware
	a.echo.Use(sentryecho.New(sentryecho.Options{
		Repanic:         true,
		WaitForDelivery: true,
	}))

	r := a.echo.Group("/api")

	r.GET("/person/", a.em.GetPersonHandler)
	r.POST("/person/", a.em.CreatePersonHandler)
	r.PUT("/person/:personId", a.em.UpdatePersonHandler)
	r.DELETE("/person/:personId", a.em.DeletePersonHandler)
	//a.echo.Use(zm.EchoIPContext(), zm.EchoSentryHubContext())
}
