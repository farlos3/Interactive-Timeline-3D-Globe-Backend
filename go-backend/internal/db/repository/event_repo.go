package repository

import (
	"context"
	"log"

	"globe/internal/db/connection"
	"globe/internal/db/models"
)

func GetEventResponses() ([]models.EventResponse, error) {
	log.Println("[DEBUG] Start querying event responses from database")
	rows, err := connection.DB.Query(context.Background(),
		`
		SELECT 
			e.event_id, e.event_name, e.date, e.lat, e.lon, e.description,
			e.country_id, c.country_name, e.event_type, e.source,
			COALESCE(ARRAY_AGG(DISTINCT t.tag_name), '{}') AS tags,
			COALESCE(ARRAY_AGG(DISTINCT ecm.cluster_id), '{}') AS clusters
		FROM event e
		LEFT JOIN country c ON e.country_id = c.country_id
		LEFT JOIN eventtag et ON e.event_id = et.event_id
		LEFT JOIN tag t ON et.tag_id = t.tag_id
		LEFT JOIN eventclustermap ecm ON e.event_id = ecm.event_id
		GROUP BY e.event_id, c.country_name
		ORDER BY e.date
	`)

	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.EventResponse

	for rows.Next() {
		var event models.EventResponse
		var tags []string
		var clusters []int

		err := rows.Scan(
			&event.EventID, &event.EventName, &event.Date,
			&event.Lat, &event.Lon, &event.Description,
			&event.CountryID, &event.CountryName, &event.EventType, &event.Source,
			&tags, &clusters,
		)
		if err != nil {
			log.Printf("[ERROR] Scanning row failed: %v", err)
			continue
		}
		log.Printf("[DEBUG] Scanned event: %+v", event)
		log.Printf("[DEBUG] Tags: %+v, Clusters: %+v", tags, clusters)
		event.Tags = tags
		event.Clusters = clusters

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Rows error: %v", err)
		return nil, err
	}

	log.Printf("[DEBUG] Total events fetched: %d", len(events))
	return events, nil
}

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
