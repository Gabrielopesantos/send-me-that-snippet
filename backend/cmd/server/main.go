package main

import (
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"time"
)

func main() {
	srv := fiber.New()
	srv.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	srv.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        10,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for")
		},
		LimitReached: func(c *fiber.Ctx) error {
			c.Status(fiber.StatusTooManyRequests)
			return c.Send([]byte("Rate limit finished"))
		},
		//Storage: myCustomStorage{}
	}))

	srv.Get("/", func(ctx *fiber.Ctx) error {
		ctx.Send([]byte("All good"))
		ctx.Status(fiber.StatusAccepted)
		return nil
	})

	srv.Post("/new", func(ctx *fiber.Ctx) error {
		var err error
		val := validator.New()
		paste := model.Paste{}
		err = ctx.BodyParser(&paste)
		if err != nil {
			log.Print(err)
			return err
		}

		err = val.Struct(&paste)
		if err != nil {
			log.Print(err)
			return err
		}

		return ctx.JSON(paste)
	})

	srv.Listen(":8080")
}
