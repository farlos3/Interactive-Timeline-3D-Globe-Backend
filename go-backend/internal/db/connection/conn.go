package connection

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

// ConnectDB ทำหน้าที่เชื่อมต่อกับฐานข้อมูล Supabase
func ConnectDB() error {
	// โหลด .env ก่อนใช้งาน
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	// ดึงค่า DATABASE_URL จาก environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return ErrDatabaseURLNotSet
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return err
	}

	err = DB.Ping(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Error messages
var (
	ErrDatabaseURLNotSet = errors.New("DATABASE_URL not set")
)
