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

func GetHierarchicalClustersHandler(c *fiber.Ctx) error {
	var query models.ClusterQuery

	// Parse request body into query struct
	if err := c.BodyParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid query parameters",
		})
	}

	// Set default max_level to 4 if not provided
	if query.MaxLevel == 0 {
		query.MaxLevel = 4
	}

	// Validate max_level
	if query.MaxLevel < 0 || query.MaxLevel > 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "MaxLevel must be between 0 and 4",
		})
	}

	// Get hierarchical clusters
	clusters, err := repository.GetHierarchicalClusters(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch hierarchical clusters",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   clusters,
	})
}
