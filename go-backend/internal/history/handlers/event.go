package handler

import (
	"globe/internal/db/models"
	"globe/internal/db/repository"

	"github.com/gofiber/fiber/v2"
)

// Response เป็นโครงสร้างมาตรฐานสำหรับการส่ง response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func GetFilteredEventsHandler(c *fiber.Ctx) error {
	var filter models.EventFilter

	// Parse request body into filter struct
	if err := c.BodyParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  "error",
			Message: "Invalid filter parameters",
			Error:   err.Error(),
		})
	}

	// Validate filter
	if filter.DateFilter != nil && filter.DateFilter.Year != nil {
		year := *filter.DateFilter.Year
		if year < 1900 || year > 2100 {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				Status:  "error",
				Message: "Invalid year range (1900-2100)",
			})
		}
	}

	// Get filtered events
	events, err := repository.GetFilteredEvents(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to fetch filtered events",
			Error:   err.Error(),
		})
	}


	return c.JSON(Response{
		Status:  "success",
		Message: "Events fetched successfully",
		Data:    events,
	})
}