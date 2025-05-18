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
			COALESCE(ARRAY_AGG(DISTINCT t.tag_name) FILTER (WHERE t.tag_name IS NOT NULL), ARRAY[]::text[]) as tags,
			COALESCE(ARRAY_AGG(DISTINCT ecm.cluster_id) FILTER (WHERE ecm.cluster_id IS NOT NULL), ARRAY[]::int[]) as clusters
		FROM event e
		LEFT JOIN eventtag et ON e.event_id = et.event_id
		LEFT JOIN tag t ON et.tag_id = t.tag_id
		LEFT JOIN eventclustermap ecm ON e.event_id = ecm.event_id
		WHERE 1=1
	`

	// 2. เพิ่มเงื่อนไข filter tags และ date
	args := []interface{}{}
	argCount := 1

	// 2.1 เพิ่มเงื่อนไข filter tags
	if filter.TagFilter != nil && len(filter.TagFilter.Tags) > 0 {
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

	// 2.2 เพิ่มเงื่อนไข filter date
	if filter.DateFilter != nil {
		if filter.DateFilter.Year != nil {
			query += fmt.Sprintf(" AND EXTRACT(YEAR FROM e.date) = $%d", argCount)
			args = append(args, *filter.DateFilter.Year)
			argCount++
		} else {
			if filter.DateFilter.StartDate != nil {
				query += fmt.Sprintf(" AND e.date >= $%d", argCount)
				args = append(args, filter.DateFilter.StartDate)
				argCount++
			}
			if filter.DateFilter.EndDate != nil {
				query += fmt.Sprintf(" AND e.date <= $%d", argCount)
				args = append(args, filter.DateFilter.EndDate)
				argCount++
			}
		}
	}

	// 3. Group by และ order by
	query += `
		GROUP BY e.event_id, e.event_name, e.date, e.lat, e.lon, e.description
		ORDER BY e.date DESC
	`

	// Debug: Print query and args
	log.Printf("[DEBUG] Args: %v", args)

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

	// Debug: Print number of results
	log.Printf("[DEBUG] Found %d events", len(events))

	// แก้ไขค่า NaN เป็น 0
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

// GetEventClusterMapping fetches individual event-cluster relationships
func GetEventClusterMapping() ([]models.EventClusterMapping, error) {
	query := `
		SELECT 
			e.event_id,
			c.cluster_id,
			c.parent_cluster_id
		FROM event e
		JOIN eventclustermap ecm ON e.event_id = ecm.event_id
		JOIN cluster c ON ecm.cluster_id = c.cluster_id
		ORDER BY e.event_id, c.cluster_id
	`

	rows, err := connection.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var mappings []models.EventClusterMapping
	for rows.Next() {
		var mapping models.EventClusterMapping
		err := rows.Scan(
			&mapping.EventID,
			&mapping.ClusterID,
			&mapping.ParentClusterID,
		)
		if err != nil {
			log.Printf("[ERROR] Scanning row failed: %v", err)
			continue
		}
		mappings = append(mappings, mapping)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Rows error: %v", err)
		return nil, err
	}

	log.Printf("[DEBUG] Found %d event-cluster mappings", len(mappings))
	return mappings, nil
}
