package main

import (
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/internal/server"
	"github.com/gabrielopesantos/smts/pkg/database"
	"log"
)

func main() {
	// Load and parse config
	cfgFile, err := config.LoadConfig("./config/config-dev.yaml")
	if err != nil {
		log.Fatalf("failed to load config file. Error: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("failed to parse config file. Error: %v", err)
	}

	// Init database connection and auto migrate
	db, err := database.NewGormDB(cfg)
	if err != nil {
		log.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&model.Paste{})
	if err != nil {
		log.Fatal("Failed to auto migrate schema")
	}

	srv := server.New(db, cfg)
	srv.Start()

	//app := fiber.New()
	//app.Use(logger.New(logger.Config{
	//	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	//}))
	//app.Use(limiter.New(limiter.Config{
	//	Next: func(c *fiber.Ctx) bool {
	//		return c.IP() == "127.0.0.1"
	//	},
	//	Max:        10,
	//	Expiration: 1 * time.Hour,
	//	KeyGenerator: func(c *fiber.Ctx) string {
	//		return c.Get("x-forwarded-for")
	//	},
	//	LimitReached: func(c *fiber.Ctx) error {
	//		c.Status(fiber.StatusTooManyRequests)
	//		return c.Send([]byte("Rate limit finished"))
	//	},
	//	//Storage: myCustomStorage{}
	//}))
	//
	//pastesGroup := app.Group("/paste")
	//
	//pastesGroup.Post("/", func(ctx *fiber.Ctx) error {
	//	var err error
	//	val := validator.New()
	//	paste := model.Paste{}
	//	err = ctx.BodyParser(&paste)
	//	if err != nil {
	//		log.Print(err)
	//		return err
	//	}
	//
	//	err = val.Struct(&paste)
	//	if err != nil {
	//		log.Print(err)
	//		return err
	//	}
	//
	//	// ?
	//	db.Create(&paste)
	//
	//	return ctx.JSON(paste)
	//})
	//
	//pastesGroup.Get("/:pId", func(ctx *fiber.Ctx) error {
	//	pId := ctx.Params("pId")
	//
	//	paste := model.Paste{}
	//	db.First(&paste, "id  = ?", pId)
	//
	//	return ctx.JSON(paste)
	//})
	//
	//mm := middleware.NewMiddlewareManager(cfg)
	//
	//pastesGroup.Delete("/:pId", mm.BasicAuthMiddleware(
	//	func(ctx *fiber.Ctx) error {
	//		pId := ctx.Params("pId")
	//
	//		paste := model.Paste{}
	//		db.Delete(&paste, "id = ?", pId)
	//
	//		return ctx.JSON(paste)
	//	}))
	//
	//// pastesGroup.Delete("/:pId", func(ctx *fiber.Ctx) error {
	//// 	pId := ctx.Params("pId")
	//
	//// 	paste := model.Paste{}
	//// 	db.Delete(&paste, "id = ?", pId)
	//
	//// 	return ctx.JSON(paste)
	//// })
	//
	//app.Listen(":8888")
}
