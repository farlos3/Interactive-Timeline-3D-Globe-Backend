package models

import "time"

type Event struct {
	EventID     int       `json:"event_id"`
	EventName   string    `json:"event_name"`
	Date        time.Time `json:"date"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Video       string    `json:"video"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
}

type EventLatLonDate struct {
	EventID int       `json:"EventID"`
	Lat     float64   `json:"Lat"`
	Lon     float64   `json:"Lon"`
	Date    time.Time `json:"Date"`
}

type EventResponse struct {
	EventID     int       `json:"event_id"`
	EventName   string    `json:"event_name"`
	Date        time.Time `json:"date"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Video       string    `json:"video"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Clusters    []int     `json:"clusters"`
}

type Cluster struct {
	ClusterID        int     `json:"cluster_id"`
	ParentClusterID  *int    `json:"parent_cluster_id"`
	CentroidLat      float64 `json:"centroid_lat"`       // จุดศูนย์กลางละติจูด
	CentroidLon      float64 `json:"centroid_lon"`       // จุดศูนย์กลางลองจิจูด
	CentroidTimeDays string  `json:"centroid_time_days"` // จุดศูนย์กลางเวลา
	Level            int     `json:"level"`
	EventIDs         []int   `json:"event_ids"`
}

type Viewport struct {
	North float64 `json:"north"` // latitude ของขอบบน
	South float64 `json:"south"` // latitude ของขอบล่าง
	East  float64 `json:"east"`  // longitude ของขอบขวา
	West  float64 `json:"west"`  // longitude ของขอบซ้าย
}

type DateFilter struct {
	StartDate *time.Time `json:"start_date"` // วันที่เริ่มต้น
	EndDate   *time.Time `json:"end_date"`   // วันที่สิ้นสุด
	Year      *int       `json:"year"`       // ปีที่ต้องการ filter
}

type ClusterQuery struct {
	Viewport    Viewport    `json:"viewport"`     // viewport ที่ user เห็น
	MaxLevel    int         `json:"max_level"`    // ระดับสูงสุดที่ต้องการ (0-2)
	TagFilter   *TagFilter  `json:"tag_filter"`   // filter ด้วย tags
	DateFilter  *DateFilter `json:"date_filter"`  // filter ด้วยวันที่
	MaxClusters *int        `json:"max_clusters"` // จำนวน clusters สูงสุดที่ต้องการ
}

type TagFilter struct {
	Tags     []string `json:"tags"`
	Operator string   `json:"operator"`
}

type EventFilter struct {
	TagFilter  *TagFilter  `json:"tag_filter"`  // ตัวเลือกสำหรับ filter tags
	DateFilter *DateFilter `json:"date_filter"` // ตัวเลือกสำหรับ filter วันที่
}
