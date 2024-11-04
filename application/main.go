package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"

	"github.com/arifnurdiansyah92/go-boilerplate/application/config"
	"github.com/arifnurdiansyah92/go-boilerplate/application/module"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func getDatabaseName(dbURL string) (string, error) {
	if dbURL == "" {
		return "", fmt.Errorf("URL is empty")
	}

	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		return "", err
	}

	dbName := strings.TrimPrefix(parsedURL.Path, "/")
	if dbName == "" {
		return "", fmt.Errorf("No database name in URL")
	}

	return dbName, nil
}

func main() {
	// Load config
	cfg, err := config.Load("/etc/cmms/config", ".")
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to load config")
	}

	// Database connection
	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}
	defer dbpool.Close()

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = &CustomValidator{validator: validator.New()}
	e.IPExtractor = echo.ExtractIPDirect()
	if cfg.WithProxy {
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
	}

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			req := c.Request()
			return strings.Contains(c.Path(), "/check") ||
				strings.Contains(c.Path(), "/me") ||
				strings.Contains(req.Method, "OPTIONS")
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// health check
	e.GET("/check", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Modules
	module.InitApps(e, dbpool, cfg)

	// Start server
	go func() {
		if err := e.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Echo listen and serve fail")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Echo shutdown error")
	}
}
