package handler

import (
	"fmt"
	"log"
	"os"
	"time"

	"globe/internal/history/service"

	"github.com/gofiber/fiber/v2"
	"github.com/go-resty/resty/v2"
)

func GetEventLatLonDateHandler(c *fiber.Ctx) error {
	// 1. ดึงข้อมูลจาก DB
	events, err := service.GetEventLatLonDate()
	if err != nil {
		log.Println("Error fetching event lat, lon, date:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch event lat, lon, date",
		})
	}

	// 2. อ่าน PY_PORT จาก env พร้อม fallback
	pyPort := os.Getenv("PY_PORT")
	pythonURL := fmt.Sprintf("http://localhost:%s/process", pyPort)

	// 3. สร้าง Resty client
	client := resty.New().SetTimeout(10 * time.Second)

	// 4. ส่ง POST ไปยัง Python server
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(events).
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