package router

import (
	"alloy/internal/shared/config"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func InitRouterWithConfig(cfg *config.Config, logger *zap.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	origins := cfg.ORIGINS
	if origins == "" {
		origins = "*"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:  origins,
		AllowMethods:  "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
		MaxAge:        300,
	}))

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		stop := time.Now()

		logger.Info("HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.OriginalURL()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", stop.Sub(start)),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.Get("User-Agent")),
		)

		return err
	})

	if strings.ToLower(cfg.APP_MODE) == "dev" && cfg.VerboseRequestLogging {
		app.Use(func(c *fiber.Ctx) error {
			contentType := c.Get("Content-Type")
			var bodyAny interface{}
			if strings.Contains(strings.ToLower(contentType), "application/json") {
				if err := json.Unmarshal(c.Body(), &bodyAny); err != nil {
					bodyAny = string(c.Body())
				}
			} else if len(c.Body()) > 0 {
				bodyAny = string(c.Body())
			} else {
				bodyAny = map[string]any{}
			}

			payload := map[string]any{
				"body":        bodyAny,
				"query":       c.Queries(),
				"params":      c.AllParams(),
				"headers":     c.GetReqHeaders(),
				"method":      c.Method(),
				"path":        c.OriginalURL(),
				"url":         c.OriginalURL(),
				"request-id":  c.Get("X-Request-Id"),
				"deployment":  c.Get("deployment"),
				"application": c.Get("application"),
			}

			b, _ := json.MarshalIndent(payload, "", "  ")
			fmt.Printf("[INFO] Request received: %s %s\n", c.Method(), c.OriginalURL())
			for _, line := range strings.Split(string(b), "\n") {
				fmt.Printf("%s\n", line)
			}

			return c.Next()
		})
	}

	return app
}

func RunWithGracefulShutdown(app *fiber.App, port string, logger *zap.Logger) error {
	go func() {
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logger.Info("🚀 🚀 Server is running",
		zap.String("url", "http://localhost:"+port),
	)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server shutdown complete.")

	return nil
}
