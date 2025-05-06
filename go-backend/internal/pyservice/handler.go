package pyservice

import (
	"log"
	"time"

	"globe/internal/db/repository"

	"github.com/gofiber/fiber/v2"
)

// Handler handles Python service related requests
type Handler struct {
	client *Client
}

// NewHandler creates a new Python service handler
func NewHandler() *Handler {
	return &Handler{
		client: NewClient(),
	}
}

// ProcessEvent handles event processing through Python service
func (h *Handler) ProcessEvent(c *fiber.Ctx) error {
	startTime := time.Now()
	log.Printf("[Python Service] Starting to process request from %s", c.IP())

	// ดึงข้อมูลจาก database
	events, err := repository.GetEventLatLonDate()
	if err != nil {
		log.Printf("[Python Service] Error fetching events: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch events from DB",
		})
	}
	log.Printf("[Python Service] Fetched %d events", len(events))

	// ส่งข้อมูลไปยัง Python service
	result, err := h.client.ProcessData(events)
	if err != nil {
		log.Printf("[Python Service] Error processing data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to process data",
			"details": err.Error(),
		})
	}

	// คำนวณเวลาที่ใช้ในการประมวลผล
	processingTime := time.Since(startTime)
	log.Printf("[Python Service] Request processed successfully in %v", processingTime)
	log.Printf("[Python Service] Result: %+v", result)

	// ส่งผลลัพธ์กลับไปยัง client
	return c.JSON(result)
}