package server

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gabrielopesantos/smts/internal/paste"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
	"log"
)

type Server struct {
	app    *fiber.App
	dbConn *gorm.DB
	cfg    *config.Config
}

func New(dbConn *gorm.DB, cfg *config.Config) *Server {
	return &Server{
		app:    fiber.New(),
		dbConn: dbConn,
		cfg:    cfg,
	}
}

func (s *Server) Start() {
	s.mapRoutes()

	addr := fmt.Sprintf("%s:%s", s.cfg.ServerConfig.Host, s.cfg.ServerConfig.Port)
	err := s.app.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) mapRoutes() {
	// This should be split
	mm := middleware.NewMiddlewareManager(s.cfg)
	if s.cfg.ServerConfig.CreateDashboardEndpoint {
		s.app.Get("/dashboard", mm.BasicAuthMiddleware(monitor.New()))
	}

	// V1
	v1Group := s.app.Group("/api/v1")

	// Paste
	pasteGroup := v1Group.Group("/paste")
	pasteHandlers := paste.NewHandlers(s.dbConn, s.cfg)
	paste.MapRoutes(pasteGroup, pasteHandlers, mm)
}
