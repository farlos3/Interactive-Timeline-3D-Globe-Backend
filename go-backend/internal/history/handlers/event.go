package handler

import (
	"fmt"
	"globe/internal/db/models"
	"globe/internal/db/repository"
	"log"

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

// GetEventClusterMappingHandler handles GET request for event cluster mapping data
func GetEventClusterMappingHandler(c *fiber.Ctx) error {
	mappings, err := repository.GetEventClusterMapping()
	if err != nil {
		log.Printf("[ERROR] Failed to get event cluster mappings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to get event cluster mappings",
			Data:    nil,
		})
	}

	return c.JSON(Response{
		Status:  "success",
		Message: "Successfully retrieved event cluster mappings",
		Data:    mappings,
	})
}

// GetEventClusterMappingCSVHandler handles GET request for event cluster mapping data in CSV format
func GetEventClusterMappingCSVHandler(c *fiber.Ctx) error {
	mappings, err := repository.GetEventClusterMapping()
	if err != nil {
		log.Printf("[ERROR] Failed to get event cluster mappings: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Status:  "error",
			Message: "Failed to get event cluster mappings",
			Data:    nil,
		})
	}

	// สร้าง CSV header
	csvData := "event_id,cluster_id,parent_cluster_id\n"

	// เพิ่มข้อมูลแต่ละแถว
	for _, mapping := range mappings {
		// แปลง parent_cluster_id เป็น string (ถ้าเป็น nil ให้เป็นค่าว่าง)
		parentClusterIDStr := ""
		if mapping.ParentClusterID != nil {
			parentClusterIDStr = fmt.Sprintf("%d", *mapping.ParentClusterID)
		}

		// สร้างแถว CSV
		csvData += fmt.Sprintf("%d,%d,%s\n",
			mapping.EventID,
			mapping.ClusterID,
			parentClusterIDStr,
		)
	}

	// ตั้งค่า header สำหรับการดาวน์โหลดไฟล์ CSV
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=event_cluster_mapping.csv")

	return c.SendString(csvData)
}