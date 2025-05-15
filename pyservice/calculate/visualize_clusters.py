import matplotlib.pyplot as plt
import numpy as np
from datetime import datetime
import folium
from folium.plugins import MarkerCluster

def visualize_clusters_on_map(events, clusters):
    """
    สร้างแผนที่แสดง events และ clusters
    """
    # สร้างแผนที่เริ่มต้นที่ centroid ของ events ทั้งหมด
    center_lat = np.mean([e['Lat'] for e in events])
    center_lon = np.mean([e['Lon'] for e in events])
    m = folium.Map(location=[center_lat, center_lon], zoom_start=10)

    # สร้าง marker cluster สำหรับ events
    marker_cluster = MarkerCluster().add_to(m)

    # เพิ่ม markers สำหรับ events
    for event in events:
        folium.Marker(
            location=[event['Lat'], event['Lon']],
            popup=f"Event ID: {event['EventID']}<br>Date: {event['Date']}",
            icon=folium.Icon(color='blue', icon='info-sign')
        ).add_to(marker_cluster)

    # เพิ่ม circles สำหรับ clusters
    colors = ['red', 'green', 'purple', 'orange', 'darkred']
    for cluster in clusters:
        color = colors[cluster['level'] % len(colors)]
        folium.Circle(
            location=[cluster['centroid_lat'], cluster['centroid_lon']],
            radius=1000,  # 1km
            popup=f"Cluster ID: {cluster['cluster_id']}<br>Level: {cluster['level']}<br>Events: {len(cluster['event_ids'])}",
            color=color,
            fill=True,
            fill_opacity=0.2
        ).add_to(m)

    # บันทึกแผนที่
    m.save('cluster_visualization.html')

def plot_cluster_hierarchy(clusters):
    """
    สร้างกราฟแสดงลำดับชั้นของ clusters
    """
    plt.figure(figsize=(12, 8))
    
    # สร้างกราฟสำหรับแต่ละ level
    for level in range(max(c['level'] for c in clusters) + 1):
        level_clusters = [c for c in clusters if c['level'] == level]
        x = [c['centroid_lon'] for c in level_clusters]
        y = [c['centroid_lat'] for c in level_clusters]
        sizes = [len(c['event_ids']) * 100 for c in level_clusters]
        
        plt.scatter(x, y, s=sizes, alpha=0.5, label=f'Level {level}')

    plt.title('Cluster Hierarchy')
    plt.xlabel('Longitude')
    plt.ylabel('Latitude')
    plt.legend()
    plt.grid(True)
    plt.savefig('cluster_hierarchy.png')
    plt.close()

def analyze_cluster_quality(events, clusters):
    """
    วิเคราะห์คุณภาพของ clusters
    """
    results = {
        'total_events': len(events),
        'total_clusters': len(clusters),
        'events_per_cluster': [],
        'cluster_levels': {},
        'spatial_spread': []
    }

    # วิเคราะห์จำนวน events ต่อ cluster
    for cluster in clusters:
        results['events_per_cluster'].append(len(cluster['event_ids']))
        results['cluster_levels'][cluster['level']] = results['cluster_levels'].get(cluster['level'], 0) + 1

        # คำนวณ spatial spread
        cluster_events = [e for e in events if e['EventID'] in cluster['event_ids']]
        if cluster_events:
            lat_spread = max(e['Lat'] for e in cluster_events) - min(e['Lat'] for e in cluster_events)
            lon_spread = max(e['Lon'] for e in cluster_events) - min(e['Lon'] for e in cluster_events)
            results['spatial_spread'].append((lat_spread, lon_spread))

    # พิมพ์ผลการวิเคราะห์
    print("\nCluster Analysis Results:")
    print(f"Total Events: {results['total_events']}")
    print(f"Total Clusters: {results['total_clusters']}")
    print(f"Average Events per Cluster: {np.mean(results['events_per_cluster']):.2f}")
    print("\nClusters per Level:")
    for level, count in sorted(results['cluster_levels'].items()):
        print(f"Level {level}: {count} clusters")
    print("\nSpatial Spread (average):")
    avg_lat_spread = np.mean([s[0] for s in results['spatial_spread']])
    avg_lon_spread = np.mean([s[1] for s in results['spatial_spread']])
    print(f"Latitude: {avg_lat_spread:.6f} degrees")
    print(f"Longitude: {avg_lon_spread:.6f} degrees")

    return results

if __name__ == "__main__":
    # ตัวอย่างการใช้งาน
    from divisive_tree import build_divisive_tree
    
    # สร้างข้อมูลทดสอบ
    test_events = [
        {
            'EventID': i,
            'Lat': 13.7563 + np.random.normal(0, 0.01),
            'Lon': 100.5018 + np.random.normal(0, 0.01),
            'Date': '2024-01-01'
        }
        for i in range(100)
    ]
    
    # สร้าง clusters
    clusters = build_divisive_tree(test_events, max_level=2)
    
    # วิเคราะห์และแสดงผล
    analyze_cluster_quality(test_events, clusters)
    visualize_clusters_on_map(test_events, clusters)
    plot_cluster_hierarchy(clusters) 