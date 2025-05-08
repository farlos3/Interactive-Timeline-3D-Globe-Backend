package models

import "time"

type Event struct {
	EventID     int       `json:"event_id"`
	EventName   string    `json:"event_name"`
	Date        time.Time `json:"date"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Description string    `json:"description"`
	CountryID   int       `json:"country_id"`
	EventType   string    `json:"event_type"`
	Source      string    `json:"source"`
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
	Description string    `json:"description"`
	CountryID   int       `json:"country_id"`
	CountryName string    `json:"country_name"`
	EventType   string    `json:"event_type"`
	Source      string    `json:"source"`
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
	GroupTag         string  `json:"group_tag"`
	BoundingBox      string  `json:"bounding_box"`
	EventIDs         []int   `json:"event_ids"`
}