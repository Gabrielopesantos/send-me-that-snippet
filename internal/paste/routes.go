package paste

import (
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(router fiber.Router, handlers *pasteHandlers, mm *middleware.Manager) {
	router.Use(mm.LoggerMiddleware())
	router.Get("/filter", mm.BasicAuthMiddleware(handlers.Filter()))
	router.Get("/:pId", handlers.Get())
	router.Post("/", handlers.Insert())
	router.Delete("/:pId", mm.BasicAuthMiddleware(handlers.Delete()))
	router.Put("/:pId", mm.BasicAuthMiddleware(handlers.Update()))
}
