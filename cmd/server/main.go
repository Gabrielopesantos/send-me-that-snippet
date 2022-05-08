package main

import (
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	pgDsn := "postgres://gabriel:gabriel@localhost:5432/main"
	db, err := gorm.Open(postgres.Open(pgDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&model.Paste{})
	if err != nil {
		log.Fatal("Failed to auto migrate schema")
	}

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

	pastesGroup := srv.Group("/paste")

	pastesGroup.Post("/", func(ctx *fiber.Ctx) error {
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

		// ?
		db.Create(&paste)

		return ctx.JSON(paste)
	})

	pastesGroup.Get("/:pId", func(ctx *fiber.Ctx) error {
		pId := ctx.Params("pId")

		paste := model.Paste{}
		db.First(&paste, "id  = ?", pId)

		return ctx.JSON(paste)
	})

	srv.Listen(":8888")
}
