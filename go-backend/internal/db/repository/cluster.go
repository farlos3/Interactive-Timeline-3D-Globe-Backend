package repository

import (
	"context"
	"log"

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
				centroid_time_days, level, group_tag, bounding_box
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (cluster_id) DO NOTHING
			`,
			cluster.ClusterID,
			cluster.ParentClusterID,
			cluster.CentroidLat,
			cluster.CentroidLon,
			cluster.CentroidTimeDays,
			cluster.Level,
			cluster.GroupTag,
			cluster.BoundingBox,
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
