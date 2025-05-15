package repository

import (
	"context"
	"log"
	"time"

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
	log.Println("[DEBUG] Start querying hierarchical clusters (recursive BBOX & date)")

	// 1. Query cluster ทั้งหมดที่ level <= max_level พร้อม min/max lat/lon, min/max date
	baseQuery := `
		SELECT 
    		c.cluster_id,
    		c.parent_cluster_id,
    		c.centroid_lat,
    		c.centroid_lon,
    		c.centroid_time_days,
    		c.level,
    				ARRAY_AGG(DISTINCT ecm.event_id) as event_ids,
    		c.min_lat, c.max_lat, c.min_lon, c.max_lon,
    		c.min_date, c.max_date
		FROM cluster c
		LEFT JOIN eventclustermap ecm ON c.cluster_id = ecm.cluster_id
		WHERE c.level <= $1
		GROUP BY c.cluster_id, c.parent_cluster_id, c.centroid_lat, c.centroid_lon, c.centroid_time_days, c.level, c.min_lat, c.max_lat, c.min_lon, c.max_lon, c.min_date, c.max_date
	`
	rows, err := connection.DB.Query(context.Background(), baseQuery, query.MaxLevel)
	if err != nil {
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var clusters []models.Cluster
	parentMap := make(map[int][]models.Cluster)
	clusterMap := make(map[int]models.Cluster)

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
			&cluster.MinLat, &cluster.MaxLat, &cluster.MinLon, &cluster.MaxLon,
			&cluster.MinDate, &cluster.MaxDate,
		)
		if err != nil {
			log.Printf("[ERROR] Scanning row failed: %v", err)
			continue
		}
		clusters = append(clusters, cluster)
		pid := 0
		if cluster.ParentClusterID != nil {
			pid = *cluster.ParentClusterID
		}
		parentMap[pid] = append(parentMap[pid], cluster)
		clusterMap[cluster.ClusterID] = cluster
	}
	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Rows error: %v", err)
		return nil, err
	}

	// Helper: ตรวจสอบว่า cluster เป็น leaf (ไม่มีลูก)
	isLeaf := func(clusterID int) bool {
		children, ok := parentMap[clusterID]
		return !ok || len(children) == 0
	}

	// 2. Recursive filter
	var result []models.Cluster
	var traverse func(parentID int)
	traverse = func(parentID int) {
		for _, c := range parentMap[parentID] {
			// เช็ค BBOX กับ viewport
			if c.MaxLat == nil || c.MinLat == nil || c.MaxLon == nil || c.MinLon == nil {
				continue // ไม่มี event ใน cluster นี้
			}
			if *c.MaxLat < query.Viewport.South || *c.MinLat > query.Viewport.North ||
				*c.MaxLon < query.Viewport.West || *c.MinLon > query.Viewport.East {
				continue
			}
			// เช็คช่วงเวลา
			if c.MaxDate == nil || c.MinDate == nil {
				continue
			}
			if query.DateFilter != nil {
				if query.DateFilter.Year != nil {
					year := *query.DateFilter.Year
					start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
					end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
					if c.MaxDate.Before(start) || c.MinDate.After(end) {
						continue
					}
				} else if query.DateFilter.StartDate != nil && query.DateFilter.EndDate != nil {
					if c.MaxDate.Before(*query.DateFilter.StartDate) || c.MinDate.After(*query.DateFilter.EndDate) {
						continue
					}
				}
			}
			// ถ้าเป็น leaf (ไม่มีลูก) → ดึง event ที่เข้าเงื่อนไข (สามารถ filter เพิ่มเติมได้ที่นี่)
			if isLeaf(c.ClusterID) {
				// (Optional) filter event รายตัวที่อยู่ใน viewport/ช่วงเวลา
				// ...
			} else {
				traverse(c.ClusterID)
			}
			result = append(result, c)
		}
	}
	traverse(0) // เริ่มจาก root (parentID = 0)

	log.Printf("[DEBUG] Total clusters after recursive filter: %d", len(result))
	return result, nil
}
