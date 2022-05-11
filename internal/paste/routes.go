package paste

import "github.com/gofiber/fiber/v2"

func MapRoutes(router fiber.Router, handlers *pasteHandlers) {
	router.Get("/:pId", handlers.Get())
	router.Post("/", handlers.Insert())
	router.Delete("/:pId", handlers.Delete())
}
