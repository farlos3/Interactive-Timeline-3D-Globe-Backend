package handler

import (
	"globe/internal/db/models"
	"globe/internal/db/repository"

	"github.com/gofiber/fiber/v2"
)

func InsertClustersHandler(c *fiber.Ctx) error {
	var clusters []models.Cluster

	// Parse JSON body เป็น slice ของ Cluster
	if err := c.BodyParser(&clusters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON",
		})
	}

	// Insert clusters and mappings
	if err := repository.InsertClustersAndMappings(clusters); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to insert clusters",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Clusters inserted successfully",
	})
}
