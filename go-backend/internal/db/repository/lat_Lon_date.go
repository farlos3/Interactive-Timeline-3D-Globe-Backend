package repository

import (
	"context"
	"log"

	"globe/internal/db/connection"
	"globe/internal/db/models"
)

func GetEventLatLonDate() ([]models.EventLatLonDate, error) {
	log.Println("[DEBUG] Start querying event lat, lon, date from database")
	rows, err := connection.DB.Query(context.Background(),
		`SELECT event_id, lat, lon, date FROM event`)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.EventLatLonDate

	for rows.Next() {
		var event models.EventLatLonDate
		err := rows.Scan(&event.EventID, &event.Lat, &event.Lon, &event.Date)
		if err != nil {
			log.Printf("[ERROR] Scanning row failed: %v", err)
			continue
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Rows error: %v", err)
		return nil, err
	}

	log.Printf("[DEBUG] Total events fetched: %d", len(events))
	return events, nil
}