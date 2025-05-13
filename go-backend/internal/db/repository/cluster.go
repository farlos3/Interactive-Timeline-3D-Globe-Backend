package repository

import (
	"context"
	"fmt"
	"log"
	"strings"

	"globe/internal/db/connection"
	"globe/internal/db/models"
)

// InsertClustersAndMappings inserts clusters and their event mappings into the database.
func InsertClustersAndMappings(clusters []models.Cluster) error {
	for _, cluster := range clusters {
		// Insert cluster
		_, err := connection.DB.Exec(context.Background(),
			`INSERT INTO cluster (
				cluster_id, parent_cluster_id, centroid_lat, centroid_lon,
				centroid_time_days, level
			) VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (cluster_id) DO NOTHING
			`,
			cluster.ClusterID,
			cluster.ParentClusterID,
			cluster.CentroidLat,
			cluster.CentroidLon,
			cluster.CentroidTimeDays,
			cluster.Level,
		)
		if err != nil {
			log.Printf("Insert cluster error: %v", err)
			return err
		}

		// Insert event-cluster mapping
		for _, eventID := range cluster.EventIDs {
			_, err := connection.DB.Exec(context.Background(),
				`INSERT INTO eventclustermap (event_id, cluster_id)
				VALUES ($1, $2)
				ON CONFLICT (event_id, cluster_id) DO NOTHING
				`, eventID, cluster.ClusterID)
			if err != nil {
				log.Printf("Insert eventclustermap error: %v", err)
				return err
			}
		}
	}
	return nil
}

// GetHierarchicalClusters ดึง clusters แบบ hierarchical ตาม viewport และ filter
func GetHierarchicalClusters(query models.ClusterQuery) ([]models.Cluster, error) {
	log.Println("[DEBUG] Start querying hierarchical clusters")

	// 1. ดึง clusters ระดับบนสุด (0-2) ที่อยู่ใน viewport
	baseQuery := `
		WITH RECURSIVE cluster_tree AS (
			-- Base case: เลือก clusters ระดับบนสุดที่อยู่ใน viewport
			SELECT 
				c.cluster_id,
				c.parent_cluster_id,
				c.centroid_lat,
				c.centroid_lon,
				c.level,
				ARRAY_AGG(DISTINCT e.event_id) as event_ids,
				ARRAY[c.cluster_id] as path,
				0 as depth
			FROM cluster c
			LEFT JOIN eventclustermap ecm ON c.cluster_id = ecm.cluster_id
			LEFT JOIN event e ON ecm.event_id = e.event_id
			WHERE c.level <= $1
			AND c.centroid_lat <= $2 AND c.centroid_lat >= $3
			AND c.centroid_lon <= $4 AND c.centroid_lon >= $5
			AND c.parent_cluster_id IS NULL
			GROUP BY c.cluster_id, c.parent_cluster_id, c.centroid_lat, c.centroid_lon, c.level

			UNION ALL

			-- Recursive case: ดึง child clusters
			SELECT 
				c.cluster_id,
				c.parent_cluster_id,
				c.centroid_lat,
				c.centroid_lon,
				c.level,
				(SELECT ARRAY_AGG(DISTINCT e.event_id)
				 FROM eventclustermap ecm
				 LEFT JOIN event e ON ecm.event_id = e.event_id
				 WHERE ecm.cluster_id = c.cluster_id) as event_ids,
				ct.path || c.cluster_id,
				ct.depth + 1
			FROM cluster c
			JOIN cluster_tree ct ON c.parent_cluster_id = ct.cluster_id
			WHERE c.level <= $1
			AND c.centroid_lat <= $2 AND c.centroid_lat >= $3
			AND c.centroid_lon <= $4 AND c.centroid_lon >= $5
		)
	`

	// 2. เพิ่มเงื่อนไข filter tags ถ้ามี
	filterConditions := ""
	args := []interface{}{
		query.MaxLevel,
		query.Viewport.North,
		query.Viewport.South,
		query.Viewport.East,
		query.Viewport.West,
	}
	argCount := 6

	// 2.1 เพิ่มเงื่อนไข filter tags
	if query.TagFilter != nil {
		if len(query.TagFilter.Tags) > 0 {
			operator := query.TagFilter.Operator
			if operator != "AND" && operator != "OR" {
				operator = "OR" // default เป็น OR
			}

			if operator == "AND" {
				for _, tag := range query.TagFilter.Tags {
					filterConditions += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM event e JOIN eventclustermap ecm ON e.event_id = ecm.event_id JOIN eventtag et ON e.event_id = et.event_id JOIN tag t ON et.tag_id = t.tag_id WHERE ecm.cluster_id = ANY(ct.event_ids) AND t.tag_name ILIKE $%d)", argCount)
					args = append(args, "%"+tag+"%")
					argCount++
				}
			} else {
				orConds := []string{}
				for _, tag := range query.TagFilter.Tags {
					orConds = append(orConds, fmt.Sprintf("t.tag_name ILIKE $%d", argCount))
					args = append(args, "%"+tag+"%")
					argCount++
				}
				filterConditions += " AND EXISTS (SELECT 1 FROM event e JOIN eventclustermap ecm ON e.event_id = ecm.event_id JOIN eventtag et ON e.event_id = et.event_id JOIN tag t ON et.tag_id = t.tag_id WHERE ecm.cluster_id = ANY(ct.event_ids) AND (" + strings.Join(orConds, " OR ") + "))"
			}
		}
	}

	// 2.2 เพิ่มเงื่อนไข filter date
	if query.DateFilter != nil {
		if query.DateFilter.Year != nil {
			filterConditions += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM event e JOIN eventclustermap ecm ON e.event_id = ecm.event_id WHERE ecm.cluster_id = ANY(ct.event_ids) AND EXTRACT(YEAR FROM e.date) = $%d)", argCount)
			args = append(args, *query.DateFilter.Year)
			argCount++
		} else {
			if query.DateFilter.StartDate != nil {
				filterConditions += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM event e JOIN eventclustermap ecm ON e.event_id = ecm.event_id WHERE ecm.cluster_id = ANY(ct.event_ids) AND e.date >= $%d)", argCount)
				args = append(args, query.DateFilter.StartDate)
				argCount++
			}
			if query.DateFilter.EndDate != nil {
				filterConditions += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM event e JOIN eventclustermap ecm ON e.event_id = ecm.event_id WHERE ecm.cluster_id = ANY(ct.event_ids) AND e.date <= $%d)", argCount)
				args = append(args, query.DateFilter.EndDate)
				argCount++
			}
		}
	}

	// 3. สร้าง query สุดท้าย
	finalQuery := baseQuery + `
		SELECT DISTINCT ON (ct.cluster_id)
			ct.cluster_id,
			ct.parent_cluster_id,
			ct.centroid_lat,
			ct.centroid_lon,
			TO_CHAR(
				TO_TIMESTAMP(AVG(EXTRACT(EPOCH FROM e.date))),
				'YYYY-MM-DD'
			) as centroid_time_days,
			ct.level,
			ct.event_ids
		FROM cluster_tree ct
		LEFT JOIN eventclustermap ecm ON ct.cluster_id = ecm.cluster_id
		LEFT JOIN event e ON ecm.event_id = e.event_id
		WHERE 1=1 ` + filterConditions + `
		GROUP BY ct.cluster_id, ct.parent_cluster_id, ct.centroid_lat, ct.centroid_lon, ct.level, ct.event_ids
		ORDER BY ct.cluster_id
	`

	if query.MaxClusters != nil {
		finalQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, *query.MaxClusters)
	}

	// 4. Execute query
	rows, err := connection.DB.Query(context.Background(), finalQuery, args...)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var clusters []models.Cluster
	for rows.Next() {
		var cluster models.Cluster
		err := rows.Scan(
			&cluster.ClusterID,
			&cluster.ParentClusterID,
			&cluster.CentroidLat,
			&cluster.CentroidLon,
			&cluster.CentroidTimeDays,
			&cluster.Level,
			&cluster.EventIDs,
		)
		if err != nil {
			log.Printf("[ERROR] Scanning row failed: %v", err)
			continue
		}
		clusters = append(clusters, cluster)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Rows error: %v", err)
		return nil, err
	}

	log.Printf("[DEBUG] Total clusters fetched: %d", len(clusters))
	return clusters, nil
}
