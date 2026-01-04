//go:build dev
// +build dev

package api

import (
	_ "github.com/AaronDevStack/CRUD_FullStackSolution/Backend/internal/services/api/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// RegisterDevRoutes registers development-only routes
func RegisterDevRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
