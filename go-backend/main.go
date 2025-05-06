package main

import (
	"log"
	"os"
	"time"

	"globe/internal/db/connection"
	"globe/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// ตั้งค่า logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("🚀 Starting Globe API Server...")

	// สร้าง Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
	})

	// เพิ่ม middleware
	app.Use(recover.New()) // จัดการ panic
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	// เชื่อมต่อฐานข้อมูล
	log.Println("📦 Connecting to database...")
	if err := connection.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// ลงทะเบียน routes ทั้งหมด
	routes.RegisterRoutes(app)

	// กำหนด port จาก environment variable หรือใช้ค่า default
	goPort := os.Getenv("GO_PORT")

	// รัน server
	log.Printf("🌐 Server is running on http://localhost:%s", goPort)
	if err := app.Listen(":" + goPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}