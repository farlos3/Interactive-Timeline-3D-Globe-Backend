package routes

import (
	"globe/internal/history/handlers"
	"globe/internal/pyservice"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes จะเชื่อม handler กับ path
func RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "🚀 Hello, Fiber + Globe API is running!",
			"time":    "now",
		})
	})

	api := app.Group("/api")
	// "Content-Type", "application/json"
	api.Post("/events-lat-lon-date", handler.GetEventLatLonDateHandler)

	// Python service routes
	pythonHandler := pyservice.NewHandler()
	api.Post("/process", pythonHandler.ProcessEvent)
}
