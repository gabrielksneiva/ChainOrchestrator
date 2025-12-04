package middleware_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gabrielksneiva/ChainOrchestrator/internal/infrastructure/logger"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_FiberError(t *testing.T) {
	log := logger.NewNopLogger()
	errorHandler := middleware.ErrorHandler(log)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "test error")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestErrorHandler_GenericError(t *testing.T) {
	log := logger.NewNopLogger()
	errorHandler := middleware.ErrorHandler(log)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return errors.New("generic error")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestErrorHandler_NoError(t *testing.T) {
	log := logger.NewNopLogger()
	errorHandler := middleware.ErrorHandler(log)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
