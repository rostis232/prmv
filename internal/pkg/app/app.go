package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/rostis232/prmv/docs"
	"github.com/rostis232/prmv/internal/handler"
	"github.com/rostis232/prmv/internal/postgres"
	"github.com/rostis232/prmv/internal/service"
	_ "github.com/swaggo/echo-swagger"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type App struct {
	Server  *echo.Echo
	Handler *handler.Handler
	Service *service.Service
}

func NewApp(pgConfig string) (*App, error) {
	var a App

	pg, err := postgres.NewPostgres(pgConfig)
	if err != nil {
		return nil, fmt.Errorf("app: failed to connect to postgres: %w", err)
	}

	err = pg.Migrate()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate postgres schema: %w", err)
	}

	a.Server = echo.New()
	a.Service = service.NewService(pg)
	a.Handler = handler.NewHandler(a.Service)
	a.Server.Use(middleware.Logger())
	a.Server.Use(middleware.Recover())

	//endpoints
	a.Server.Any("/", a.Handler.Home)
	a.Server.POST("/posts", a.Handler.AddPost)
	a.Server.GET("/posts", a.Handler.GetAllPosts)
	a.Server.PUT("/posts/:id", a.Handler.UpdatePost)
	a.Server.GET("/posts/:id", a.Handler.GetPost)
	a.Server.DELETE("/posts/:id", a.Handler.DeletePost)
	//swagger
	a.Server.GET("/swagger/*", echoSwagger.WrapHandler)

	return &a, nil
}

func (a *App) Run(port string) error {
	log.Info("app starting")
	return a.Server.Start(":" + port)
}
