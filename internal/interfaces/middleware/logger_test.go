package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestLoggerMiddleware_Success(t *testing.T) {
	log := logger.NewNopLogger()

	app := fiber.New()
	app.Use(middleware.LoggerMiddleware(log))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLoggerMiddleware_Error(t *testing.T) {
	log := logger.NewNopLogger()

	app := fiber.New()
	app.Use(middleware.LoggerMiddleware(log))

	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(500, "internal error")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestLoggerMiddleware_POST(t *testing.T) {
	log := logger.NewNopLogger()

	app := fiber.New()
	app.Use(middleware.LoggerMiddleware(log))

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(201)
	})

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestLoggerMiddleware_DifferentPaths(t *testing.T) {
	log := logger.NewNopLogger()

	app := fiber.New()
	app.Use(middleware.LoggerMiddleware(log))

	paths := []string{"/health", "/transaction", "/walletbalance"}

	for _, path := range paths {
		app.Get(path, func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})

		req := httptest.NewRequest("GET", path, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}
}
