//go:build !dev
// +build !dev

package api

import "github.com/gofiber/fiber/v2"

// RegisterDevRoutes is a no-op in production
func RegisterDevRoutes(app *fiber.App) {
	// No development routes in production
}
