package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	srv := fiber.New()

	srv.Get("/", func(ctx *fiber.Ctx) error {
		ctx.Send([]byte("All good"))
		ctx.Status(fiber.StatusAccepted)
		return nil
	})

	srv.Listen(":8080")
}
