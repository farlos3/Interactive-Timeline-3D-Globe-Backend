package pyservice

import (
    "encoding/json"
    "fmt"
    "github.com/gofiber/fiber/v2"
)

// กำหนด URL ของ Python service
const (
    baseURL = "http://localhost:8000"
)

// Client struct เก็บข้อมูลที่จำเป็นสำหรับการเชื่อมต่อ
type Client struct {
    baseURL string    // URL ของ Python service
    app     *fiber.App  // Fiber app instance
}

// สร้าง Client ใหม่
func NewClient() *Client {
    return &Client{
        baseURL: baseURL,
        app:     fiber.New(),
    }
}

// ProcessData ส่งข้อมูลไปยัง Python service
func (c *Client) ProcessData(data interface{}) (map[string]interface{}, error) {
    // 1. แปลงข้อมูลเป็น JSON
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("error marshaling data: %v", err)
    }

    // 2. สร้าง HTTP agent สำหรับส่ง request
    agent := fiber.AcquireAgent()
    defer fiber.ReleaseAgent(agent)  // ปล่อย agent เมื่อเสร็จ

    // 3. ตั้งค่า HTTP request
    req := agent.Request()
    req.Header.SetMethod(fiber.MethodPost)  // ใช้ HTTP POST
    req.SetRequestURI(fmt.Sprintf("%s/process", c.baseURL))  // ตั้งค่า URL
    req.Header.SetContentType("application/json")  // กำหนด content type เป็น JSON
    req.SetBody(jsonData)  // ใส่ข้อมูลที่จะส่ง

    // 4. ส่ง request
    if err := agent.Parse(); err != nil {
        return nil, fmt.Errorf("error parsing request: %v", err)
    }

    // 5. รับ response
    code, body, errs := agent.Bytes()
    if len(errs) > 0 {
        return nil, fmt.Errorf("error sending request: %v", errs[0])
    }

    // 6. ตรวจสอบ status code
    if code != fiber.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", code)
    }

    // 7. แปลง response เป็น map
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("error decoding response: %v", err)
    }

    // 8. ส่งผลลัพธ์กลับ
    return result, nil
}