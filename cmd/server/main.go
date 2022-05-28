package main

import (
	"context"
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/internal/server"
	"github.com/gabrielopesantos/smts/pkg/database"
	"github.com/gabrielopesantos/smts/pkg/logger"
	sntry "github.com/gabrielopesantos/smts/pkg/sentry"
	provider "github.com/gabrielopesantos/smts/pkg/tracer_provider"
	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Seed to generation of strings
	rand.Seed(time.Now().UnixNano())

	///

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := provider.InitProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer("test-tracer")

	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	ctx, span := tracer.Start(
		ctx,
		"CollectorExporter-Example",
		trace.WithAttributes(commonAttrs...))
	defer span.End()
	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)

		<-time.After(time.Second)
		iSpan.End()
	}

	///

	// Sentry init
	err = sntry.InitClient()
	if err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Load and parse config
	cfgFile, err := config.LoadConfig("./config/config-dev.yaml")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("failed to load config file. Error: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("failed to parse config file. Error: %v", err)
	}

	// Init database connection and auto migrate
	db, err := database.NewGormDB(cfg)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal("failed to connect database")
	}
	err = db.AutoMigrate(&model.Paste{})
	if err != nil {
		sentry.CaptureException(err)
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
