package http

import (
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/handlers"
	"github.com/gabrielksneiva/ChainOrchestrator/internal/interfaces/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

// Router configuração de rotas do Fiber
type Router struct {
	app                *fiber.App
	transactionHandler *handlers.TransactionHandler
	logger             *zap.Logger
}

// NewRouter cria um novo router
func NewRouter(
	transactionHandler *handlers.TransactionHandler,
	logger *zap.Logger,
) *Router {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(logger),
		AppName:      "ChainOrchestrator",
	})

	return &Router{
		app:                app,
		transactionHandler: transactionHandler,
		logger:             logger,
	}
}

// Setup configura todas as rotas
func (r *Router) Setup() {
	// Middlewares globais
	r.app.Use(recover.New())
	r.app.Use(cors.New())
	r.app.Use(middleware.LoggerMiddleware(r.logger))

	// Health check
	r.app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "ChainOrchestrator",
		})
	})

	// POST /transaction - Publica transação no SNS
	r.app.Post("/transaction", r.transactionHandler.PostTransaction)
}

// GetApp retorna a aplicação Fiber
func (r *Router) GetApp() *fiber.App {
	return r.app
}
