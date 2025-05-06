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
    EventID int
    Lat     float64
    Lon     float64
    Date    time.Time
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

