package service

import (
    "log"
    "globe/internal/db/models"
    "globe/internal/db/repository"
)

func GetEventLatLonDate() ([]models.EventLatLonDate, error) {
    events, err := repository.GetEventLatLonDate()
    if err != nil {
        log.Println("Error fetching event lat, lon, date:", err)
        return nil, err
    }
    return events, nil
}