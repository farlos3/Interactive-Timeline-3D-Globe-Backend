package models

import (
	"time"
)

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
	CentroidLat      float64 `json:"centroid_lat"`
	CentroidLon      float64 `json:"centroid_lon"`
	CentroidTimeDays string  `json:"centroid_time_days"`
	Level            int     `json:"level"`
	EventIDs         []int   `json:"event_ids"`

	Events  []EventResponse `json:"events"`
	MinLat  *float64        `json:"min_lat"`
	MaxLat  *float64        `json:"max_lat"`
	MinLon  *float64        `json:"min_lon"`
	MaxLon  *float64        `json:"max_lon"`
	MinDate *time.Time      `json:"min_date"`
	MaxDate *time.Time      `json:"max_date"`
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
	MaxLevel    int         `json:"max_level"`    // ระดับสูงสุดที่ต้องการ (0-4)
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

type EventFull struct {
	EventID     int      `json:"event_id"`
	EventName   string   `json:"event_name"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	Image       string   `json:"image"`
	Video       string   `json:"video"`
	Tags        []string `json:"tags"`
	Clusters    []int    `json:"clusters"`
}

type ClusterResponse struct {
	ClusterID        int         `json:"cluster_id"`
	ParentClusterID  *int        `json:"parent_cluster_id"`
	CentroidLat      float64     `json:"centroid_lat"`
	CentroidLon      float64     `json:"centroid_lon"`
	CentroidTimeDays float64     `json:"centroid_time_days"`
	Level            int         `json:"level"`
	MinDate          time.Time   `json:"min_date"`
	MaxDate          time.Time   `json:"max_date"`
	MinLat           float64     `json:"min_lat"`
	MaxLat           float64     `json:"max_lat"`
	MinLon           float64     `json:"min_lon"`
	MaxLon           float64     `json:"max_lon"`
	Events           []EventFull `json:"events"` // สำคัญ!
}

type EventClusterMapping struct {
	EventID         int  `json:"event_id"`
	ClusterID       int  `json:"cluster_id"`
	ParentClusterID *int `json:"parent_cluster_id"`
}
