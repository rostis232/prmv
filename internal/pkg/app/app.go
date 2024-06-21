package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rostis232/prmv/internal/handler"
	"github.com/rostis232/prmv/internal/postgres"
	"github.com/rostis232/prmv/internal/service"
)

type App struct {
	Server  *echo.Echo
	Handler *handler.Handler
	Service *service.Service
}

func NewApp(pgConfig string) (*App, error) {
	a := App{}
	pg, err := postgres.NewPostgres(pgConfig)
	if err != nil {
		return nil, err
	}
	a.Server = echo.New()
	a.Service = service.NewService(pg)
	a.Handler = handler.NewHandler(a.Service)
	a.Server.Use(middleware.Logger())
	a.Server.Use(middleware.Recover())
	a.Server.Static("/static", "./static")
	//endpoints
	a.Server.POST("/posts", a.Handler.AddPost)
	a.Server.GET("/posts", a.Handler.GetAllPosts)
	a.Server.PUT("/posts/:id", a.Handler.UpdatePost)
	a.Server.GET("/posts/:id", a.Handler.GetPost)
	a.Server.DELETE("/posts/:id", a.Handler.DeletePost)
	return &a, nil
}

func (a *App) Run(port string) error {
	return a.Server.Start(":" + port)
}
