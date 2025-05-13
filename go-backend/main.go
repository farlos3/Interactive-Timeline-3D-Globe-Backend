package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"globe/internal/db/connection"
	"globe/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Config logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("🚀 Starting Globe API Server...")

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  5 * time.Second,
		AppName:      "Globe API",
	})

	// middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path} | ${ip} | ${headers}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
		Output:     os.Stdout,
	}))

	// ตั้งค่า CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // รองรับทั้ง localhost, IP และ local IP
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowCredentials: true,
		MaxAge:           300, // ระยะเวลาที่ browser เก็บ cache preflight response (วินาที)
	}))

	log.Println("📦 Connecting to database...")
	if err := connection.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	routes.RegisterRoutes(app)

	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("❌ Failed to load .env file from %s: %v", envPath, err)
	}

	goPort := os.Getenv("GO_PORT")
	
	log.Printf("🌐 Server is running on http://localhost:%s", goPort)
	if err := app.Listen(":" + goPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
