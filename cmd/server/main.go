package main

import (
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/internal/server"
	"github.com/gabrielopesantos/smts/pkg/database"
	"github.com/gabrielopesantos/smts/pkg/logger"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Seed to generation of strings
	rand.Seed(time.Now().UnixNano())
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
		log.Fatalf("failed to auto migrate schema. Error: %v", err)
	}

	// Init logger
	loggr := logger.New(os.Stdout, logger.Info, logger.Json, false)

	// Create middleware manager ?
	middlewareManager := middleware.NewMiddlewareManager(cfg)

	// Init and start server
	srv := server.New(db, loggr, middlewareManager, cfg)
	srv.Start()
}
