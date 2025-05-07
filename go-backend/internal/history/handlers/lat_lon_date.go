package handler

import (
	"fmt"
	"log"
	"os"
	"time"
	"encoding/json"

	"globe/internal/history/service"
	"globe/internal/db/models"
    "globe/internal/db/repository"

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

    requestBody := fiber.Map{
        "events": events,
    }

    client := resty.New().SetTimeout(10 * time.Second)

    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetBody(requestBody).
        Post(pythonURL)

    if err != nil {
        log.Println("Error sending request to Python:", err)
        return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
            "error": "Failed to send data to Python service",
        })
    }

    // 1. Decode response body เป็น struct
    var pyResp struct {
        Status string `json:"status"`
        Data struct {
            Clusters []models.Cluster `json:"clusters"`
            // ... field อื่นๆ ถ้าต้องการ
        } `json:"data"`
    }

    if err := json.Unmarshal(resp.Body(), &pyResp); err != nil {
        log.Println("Error decoding Python response:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Invalid response from Python service",
        })
    }

    // 2. Insert clusters ลง DB
    if err := repository.InsertClustersAndMappings(pyResp.Data.Clusters); err != nil {
        log.Println("Error inserting clusters:", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to insert clusters",
        })
    }

    // 3. ส่ง response กลับ client (หรือจะส่ง pyResp กลับไปเลยก็ได้)
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Clusters inserted successfully",
        "clusters": pyResp.Data.Clusters,
    })
}