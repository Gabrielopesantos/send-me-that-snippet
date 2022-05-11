package paste

import (
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func MapRoutes(router fiber.Router, handlers *pasteHandlers, mm *middleware.MiddlewareManager) {
	router.Use(logger.New()) // Make a better logger
	router.Get("/:pId", handlers.Get())
	router.Post("/", handlers.Insert())
	router.Delete("/:pId", mm.BasicAuthMiddleware(handlers.Delete()))
}
