package repository

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"

	"globe/internal/db/connection"
	"globe/internal/db/models"
)

func GetFilteredEvents(filter models.EventFilter) ([]models.EventResponse, error) {
	// 1. สร้าง base query
	query := `
		SELECT DISTINCT
			e.event_id,
			e.event_name,
			e.date,
			e.lat,
			e.lon,
			e.description,
			ARRAY_AGG(DISTINCT t.tag_name) as tags,
			ARRAY_AGG(DISTINCT ecm.cluster_id) as clusters
		FROM event e
		LEFT JOIN eventtag et ON e.event_id = et.event_id
		LEFT JOIN tag t ON et.tag_id = t.tag_id
		LEFT JOIN eventclustermap ecm ON e.event_id = ecm.event_id
		WHERE 1=1
	`

	// 2. เพิ่มเงื่อนไข filter tags
	args := []interface{}{}
	argCount := 1

	if filter.TagFilter != nil {
		if len(filter.TagFilter.Tags) > 0 {
			operator := filter.TagFilter.Operator
			if operator != "AND" && operator != "OR" {
				operator = "OR" // default เป็น OR
			}

			if operator == "AND" {
				for _, tag := range filter.TagFilter.Tags {
					query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM eventtag et2 JOIN tag t2 ON et2.tag_id = t2.tag_id WHERE et2.event_id = e.event_id AND t2.tag_name ILIKE $%d)", argCount)
					args = append(args, "%"+tag+"%")
					argCount++
				}
			} else {
				orConds := []string{}
				for _, tag := range filter.TagFilter.Tags {
					orConds = append(orConds, fmt.Sprintf("t2.tag_name ILIKE $%d", argCount))
					args = append(args, "%"+tag+"%")
					argCount++
				}
				query += " AND EXISTS (SELECT 1 FROM eventtag et2 JOIN tag t2 ON et2.tag_id = t2.tag_id WHERE et2.event_id = e.event_id AND (" + strings.Join(orConds, " OR ") + "))"
			}
		}
	}

	// 3. Group by และ order by
	query += `
		GROUP BY e.event_id, e.event_name, e.date, e.lat, e.lon, e.description
		ORDER BY e.date DESC
	`

	// 4. Execute query
	rows, err := connection.DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.EventResponse
	for rows.Next() {
		var event models.EventResponse
		err := rows.Scan(
			&event.EventID,
			&event.EventName,
			&event.Date,
			&event.Lat,
			&event.Lon,
			&event.Description,
			&event.Tags,
			&event.Clusters,
		)
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

	for i := range events {
		if math.IsNaN(events[i].Lat) {
			events[i].Lat = 0
		}
		if math.IsNaN(events[i].Lon) {
			events[i].Lon = 0
		}
	}

	return events, nil
}
