package app

import (
	"context"
	"emsrv/pkg/db"
	"emsrv/pkg/em"
	"emsrv/pkg/embedlog"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"time"
)

type Config struct {
	Database *pg.Options
	Server   struct {
		Host    string
		Port    int
		IsDevel bool
	}
	OpenApi struct {
		AgeUrl         string
		GenderUrl      string
		NationalityUtl string
	}
}

type App struct {
	embedlog.Logger
	appName  string
	cfg      Config
	db       db.DB
	dbc      *pg.DB
	echo     *echo.Echo
	duration string
	em       *em.EmService
}

func New(appName string, verbose bool, cfg Config, db db.DB, dbc *pg.DB) *App {
	a := &App{
		appName: appName,
		cfg:     cfg,
		db:      db,
		dbc:     dbc,
		echo:    echo.New(),
	}
	a.SetStdLoggers(verbose)
	a.echo.HideBanner = true
	a.echo.HidePort = true
	a.echo.IPExtractor = echo.ExtractIPFromRealIPHeader()
	a.em = em.NewEmService(a.Logger, a.db, em.Config{a.cfg.OpenApi.AgeUrl, a.cfg.OpenApi.GenderUrl, a.cfg.OpenApi.NationalityUtl})
	return a
}

// Run is a function that runs application.
func (a *App) Run() error {
	a.registerHandlers()
	return a.runHTTPServer(a.cfg.Server.Host, a.cfg.Server.Port)
}

// Shutdown is a function that gracefully stops HTTP server.
func (a *App) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := a.echo.Shutdown(ctx); err != nil {
		a.Errorf("shutting down server err=%q", err)
	}
}
