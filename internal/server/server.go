package server

import (
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gabrielopesantos/smts/internal/paste"
	"github.com/gabrielopesantos/smts/pkg/logger"
	"github.com/gabrielopesantos/smts/pkg/sentryfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
	"log"
	"time"
)

type Server struct {
	app    *fiber.App
	dbConn *gorm.DB
	logger *logger.Logger
	mm     *middleware.Manager
	cfg    *config.Config
}

func New(dbConn *gorm.DB, logger *logger.Logger, mm *middleware.Manager, cfg *config.Config) *Server {
	return &Server{
		app:    fiber.New(),
		dbConn: dbConn,
		logger: logger,
		mm:     mm,
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

func (s *Server) setGlobalMiddleware() {
	s.app.Use(csrf.New(), cors.New()) // ?
	s.app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        200,
		Expiration: 30 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	if s.cfg.ServerConfig.CreateDashboardEndpoint {
		s.app.Get("/dashboard", s.mm.BasicAuthMiddleware(monitor.New()))
	}

}

func (s *Server) mapRoutes() {
	s.app.Use(sentryfiber.New(sentryfiber.Options{}))
	// V1
	v1Group := s.app.Group("/api/v1")

	// Pastes
	go func() {
		paste.StartExpirePastesProcess(s.cfg.ServerConfig)
	}()
	pasteGroup := v1Group.Group("/pastes")
	pasteHandlers := paste.NewHandlers(s.dbConn, s.logger, s.cfg)
	paste.MapRoutes(pasteGroup, pasteHandlers, s.mm)
}
