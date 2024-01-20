package server

import (
	"example.com/m/config"
	"example.com/m/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type server struct {
	config *config.Config
}

type Server interface {
	Run() error
}

func NewServer(config *config.Config) Server {
	return &server{
		config: config,
	}
}

func (s *server) Run() error {
	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/")

	routes := services.NewServices(api, s.config)
	routes.RegisterRouters()

	return app.Listen(fmt.Sprintf("%v:%v", s.config.Host, s.config.Port))
}
