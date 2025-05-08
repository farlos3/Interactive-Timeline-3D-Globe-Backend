from calculate.divisive_tree import divisive_custom_tree
from calculate.visualize_clusters import create_cluster_visualization, save_cluster_visualization
import numpy as np
from datetime import datetime, timedelta

# สร้างข้อมูลตัวอย่างจำนวนมาก
def generate_sample_events():
    # กำหนดจุดศูนย์กลางของแต่ละภูมิภาค
    regions = {
        'กรุงเทพฯ': (13.7563, 100.5018),
        'เชียงใหม่': (18.7961, 98.9792),
        'ขอนแก่น': (16.4321, 102.8236),
        'สงขลา': (7.0084, 100.4767),
        'ภูเก็ต': (7.8804, 98.3923),
        'อุบลราชธานี': (15.2287, 104.8594),
        'นครราชสีมา': (14.9798, 102.0978),
        'สุราษฎร์ธานี': (9.1397, 99.3307)
    }
    
    sample_pairs = []
    start_date = datetime(2024, 1, 1)
    
    # สร้าง events สำหรับแต่ละภูมิภาค
    for region, (center_lat, center_lon) in regions.items():
        # สร้าง 20-30 events สำหรับแต่ละภูมิภาค
        num_events = np.random.randint(20, 31)
        
        for i in range(num_events):
            # สุ่มตำแหน่งรอบๆ จุดศูนย์กลาง (รัศมีประมาณ 0.1 องศา)
            lat = center_lat + np.random.normal(0, 0.05)
            lon = center_lon + np.random.normal(0, 0.05)
            
            # สุ่มวันที่ในช่วง 30 วัน
            days_offset = np.random.randint(0, 30)
            event_date = start_date + timedelta(days=days_offset)
            
            # สร้าง vector และ event
            vector = np.array([lat, lon, days_offset])
            event = {
                'lat': lat,
                'lon': lon,
                'date': event_date.strftime('%Y-%m-%d'),
                'region': region
            }
            
            sample_pairs.append((vector, event))
    
    return sample_pairs

# สร้างข้อมูลตัวอย่าง
sample_pairs = generate_sample_events()

# สร้าง cluster tree
min_cluster_size = 10  # ปรับขนาดขั้นต่ำของ cluster
cluster_tree = divisive_custom_tree(sample_pairs, min_cluster_size)

# สร้าง visualization
fig = create_cluster_visualization(cluster_tree)
fig.show()  # แสดงใน Jupyter notebook
# หรือ
save_cluster_visualization(cluster_tree, "cluster_visualization.html")  # บันทึกเป็นไฟล์ HTML 