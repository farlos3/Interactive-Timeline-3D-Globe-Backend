package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"globe/internal/history/handlers"
	"globe/internal/history/service"
	"globe/internal/pyservice"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes จะเชื่อม handler กับ path
func RegisterRoutes(app *fiber.App) {
	// Root route
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

// RegisterTestRoutes registers test routes with inline handler
func RegisterTestRoutes(app *fiber.App) {
	app.Get("/api/test-events", func(c *fiber.Ctx) error {
		// Get real data from database
		events, err := service.GetEventLatLonDate()
		if err != nil {
			log.Println("Error fetching event lat, lon, date:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch event lat, lon, date",
			})
		}

		// Convert to JSON
		jsonData, err := json.Marshal(events)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to marshal data",
			})
		}

		// Send data to Python API
		resp, err := http.Post("http://localhost:8000/process", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error sending data to Python API: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to send data to Python API",
			})
		}
		defer resp.Body.Close()

		// Read response from Python API
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode Python API response",
			})
		}

		return c.JSON(fiber.Map{
			"message":  "Test successful",
			"data":     events,
			"response": response,
		})
	})
}