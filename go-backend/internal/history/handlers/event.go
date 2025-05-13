package handler

import (
	"globe/internal/db/models"
	"globe/internal/db/repository"

	"github.com/gofiber/fiber/v2"
)

func GetFilteredEventsHandler(c *fiber.Ctx) error {
	var filter models.EventFilter

	// Parse request body into filter struct
	if err := c.BodyParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid filter parameters",
		})
	}

	// Get filtered events
	events, err := repository.GetFilteredEvents(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch filtered events",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   events,
	})
}
