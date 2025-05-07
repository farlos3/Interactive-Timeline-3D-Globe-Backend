package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	"globe/internal/history/service"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func GetEventLatLonDateHandler(c *fiber.Ctx) error {
	events, err := service.GetEventLatLonDate()
	if err != nil {
		log.Println("Error fetching event lat, lon, date:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch event lat, lon, date",
		})
	}

	pyPort := os.Getenv("PY_PORT")
	pythonURL := fmt.Sprintf("http://localhost:%s/process", pyPort)

	// สร้าง request body ในรูปแบบที่ Python ต้องการ
	requestBody := fiber.Map{
		"events": events,
	}

	client := resty.New().SetTimeout(10 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody). // ส่ง requestBody ที่มี events
		Post(pythonURL)

	if err != nil {
		log.Println("Error sending request to Python:", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Failed to send data to Python service",
		})
	}

	c.Set("Content-Type", "application/json")
	return c.Status(resp.StatusCode()).Send(resp.Body())
}
