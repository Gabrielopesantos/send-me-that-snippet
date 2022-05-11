package server

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/paste"
	"github.com/gofiber/fiber/v2"
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
	v1Group := s.app.Group("/api/v1")

	// Paste
	pasteGroup := v1Group.Group("/paste")
	pasteHandlers := paste.NewHandlers(s.dbConn, s.cfg)
	paste.MapRoutes(pasteGroup, pasteHandlers)
}
