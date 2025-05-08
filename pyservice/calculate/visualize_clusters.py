import plotly.graph_objects as go
import numpy as np
from datetime import datetime
from .divisive_tree import ClusterNode, get_leaf_clusters_with_level

def create_cluster_visualization(node, min_date=None):
    """
    สร้าง visualization ของ clusters บน globe พร้อม bounding box
    Args:
        node: ClusterNode ที่ต้องการ visualize
        min_date: วันที่เริ่มต้นสำหรับการคำนวณ centroid_time_days
    Returns:
        plotly figure object ที่สามารถแสดงผลได้
    """
    # สร้าง figure
    fig = go.Figure()
    
    # เก็บข้อมูลของแต่ละ level
    level_colors = {}  # เก็บสีสำหรับแต่ละ level
    level_nodes = {}   # เก็บ nodes ในแต่ละ level
    
    # รวบรวม nodes ตาม level
    for node, level in get_leaf_clusters_with_level(node):
        if level not in level_nodes:
            level_nodes[level] = []
        level_nodes[level].append(node)
    
    # สร้างสีสำหรับแต่ละ level
    for level in level_nodes.keys():
        # สร้างสีแบบ random แต่ให้แต่ละ level มีโทนสีใกล้เคียงกัน
        base_hue = np.random.random()  # สุ่ม hue หลัก
        level_colors[level] = [f'hsl({base_hue * 360}, 70%, {50 + i * 10}%)' 
                             for i in range(len(level_nodes[level]))]
    
    # เพิ่ม markers สำหรับแต่ละ cluster
    for level, nodes in level_nodes.items():
        for i, node in enumerate(nodes):
            # คำนวณ centroid และ bounding box
            lats = []
            lons = []
            for _, event in node.pairs:
                lat = event.get('lat', event.get('Lat'))
                lon = event.get('lon', event.get('Lon'))
                if lat is not None and lon is not None:
                    try:
                        lats.append(float(lat))
                        lons.append(float(lon))
                    except Exception:
                        continue
            
            if not lats or not lons:
                continue
                
            centroid_lat = np.mean(lats)
            centroid_lon = np.mean(lons)
            
            # คำนวณ bounding box
            min_lat, max_lat = min(lats), max(lats)
            min_lon, max_lon = min(lons), max(lons)
            
            # สร้าง bounding box coordinates
            box_lons = [min_lon, min_lon, max_lon, max_lon, min_lon]
            box_lats = [min_lat, max_lat, max_lat, min_lat, min_lat]
            
            # เพิ่ม bounding box (เส้นหนาและโปร่งแสง)
            fig.add_trace(go.Scattergeo(
                lon=box_lons,
                lat=box_lats,
                mode='lines',
                line=dict(
                    width=3,
                    color=level_colors[level][i],
                ),
                fill='toself',
                fillcolor=level_colors[level][i].replace(')', ', 0.1)'),
                name=f"Bounding Box - {node.label}",
                showlegend=True
            ))
            
            # เพิ่ม marker สำหรับ centroid
            fig.add_trace(go.Scattergeo(
                lon=[centroid_lon],
                lat=[centroid_lat],
                text=[f"Cluster {node.label}<br>Size: {len(node.pairs)}<br>Bounding Box: {min_lat:.4f}, {min_lon:.4f} to {max_lat:.4f}, {max_lon:.4f}"],
                marker=dict(
                    size=12,
                    color=level_colors[level][i],
                    symbol='circle',
                    line=dict(width=2, color='white')
                ),
                name=f"Level {level} - {node.label}"
            ))
            
            # เพิ่ม markers สำหรับ events ใน cluster
            event_lats = []
            event_lons = []
            for _, event in node.pairs:
                lat = event.get('lat', event.get('Lat'))
                lon = event.get('lon', event.get('Lon'))
                if lat is not None and lon is not None:
                    try:
                        event_lats.append(float(lat))
                        event_lons.append(float(lon))
                    except Exception:
                        continue
            
            if event_lats and event_lons:
                fig.add_trace(go.Scattergeo(
                    lon=event_lons,
                    lat=event_lats,
                    text=[f"Event in {node.label}" for _ in range(len(event_lats))],
                    marker=dict(
                        size=6,
                        color=level_colors[level][i],
                        symbol='circle',
                        opacity=0.7,
                        line=dict(width=1, color='white')
                    ),
                    showlegend=False
                ))
    
    # ปรับแต่ง layout
    fig.update_layout(
        title='Hierarchical Clustering Visualization with Bounding Boxes',
        geo=dict(
            projection_type='orthographic',
            showland=True,
            showcoastlines=True,
            showocean=True,
            oceancolor='rgb(204, 229, 255)',
            landcolor='rgb(243, 243, 243)',
            coastlinecolor='rgb(128, 128, 128)',
            showcountries=True,
            countrycolor='rgb(128, 128, 128)',
            showframe=False
        ),
        height=800,
        width=800,
        showlegend=True,
        legend=dict(
            yanchor="top",
            y=0.99,
            xanchor="left",
            x=0.01
        )
    )
    
    return fig

def save_cluster_visualization(node, output_path, min_date=None):
    """
    บันทึก visualization เป็น HTML file
    Args:
        node: ClusterNode ที่ต้องการ visualize
        output_path: path สำหรับบันทึกไฟล์ HTML
        min_date: วันที่เริ่มต้นสำหรับการคำนวณ centroid_time_days
    """
    fig = create_cluster_visualization(node, min_date)
    fig.write_html(output_path)