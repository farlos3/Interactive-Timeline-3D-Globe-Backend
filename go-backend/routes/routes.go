package routes

import (
	handler "globe/internal/history/handlers"
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
	api.Post("/events-lat-lon-date", handler.GetEventLatLonDateHandler)
	api.Post("/insert-clusters", handler.InsertClustersHandler)
	api.Post("/events/filter", handler.GetFilteredEventsHandler)
	api.Post("/clusters/hierarchical", handler.GetHierarchicalClustersHandler)
	api.Get("/events/cluster-mapping", handler.GetEventClusterMappingHandler)
	api.Get("/events/cluster-mapping/csv", handler.GetEventClusterMappingCSVHandler)

	// Python service routes
	pythonHandler := pyservice.NewHandler()
	api.Post("/process", pythonHandler.ProcessEvent)
}
