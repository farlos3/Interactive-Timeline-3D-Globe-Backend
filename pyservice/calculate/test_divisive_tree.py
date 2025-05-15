import unittest
import numpy as np
from .divisive_tree import calculate_centroid, build_divisive_tree

class TestDivisiveTree(unittest.TestCase):
    def setUp(self):
        # สร้างข้อมูลทดสอบ
        self.test_events = [
            {
                'EventID': 1,
                'Lat': 13.7563,
                'Lon': 100.5018,
                'Date': '2024-01-01'
            },
            {
                'EventID': 2,
                'Lat': 13.7564,
                'Lon': 100.5019,
                'Date': '2024-01-01'
            },
            {
                'EventID': 3,
                'Lat': 13.7565,
                'Lon': 100.5020,
                'Date': '2024-01-01'
            }
        ]

    def test_calculate_centroid(self):
        # ทดสอบการคำนวณ centroid
        centroid_lat, centroid_lon = calculate_centroid(self.test_events)
        
        # ตรวจสอบว่า centroid อยู่ระหว่างจุดข้อมูล
        self.assertGreater(centroid_lat, min(e['Lat'] for e in self.test_events))
        self.assertLess(centroid_lat, max(e['Lat'] for e in self.test_events))
        self.assertGreater(centroid_lon, min(e['Lon'] for e in self.test_events))
        self.assertLess(centroid_lon, max(e['Lon'] for e in self.test_events))

    def test_build_divisive_tree(self):
        # ทดสอบการสร้าง divisive tree
        clusters = build_divisive_tree(self.test_events, max_level=2)
        
        # ตรวจสอบโครงสร้างของ clusters
        self.assertIsInstance(clusters, list)
        self.assertGreater(len(clusters), 0)
        
        # ตรวจสอบว่าแต่ละ cluster มีข้อมูลครบ
        for cluster in clusters:
            self.assertIn('cluster_id', cluster)
            self.assertIn('centroid_lat', cluster)
            self.assertIn('centroid_lon', cluster)
            self.assertIn('level', cluster)
            self.assertIn('event_ids', cluster)
            
            # ตรวจสอบว่า centroid อยู่ในช่วงที่ถูกต้อง
            self.assertGreater(cluster['centroid_lat'], min(e['Lat'] for e in self.test_events))
            self.assertLess(cluster['centroid_lat'], max(e['Lat'] for e in self.test_events))
            self.assertGreater(cluster['centroid_lon'], min(e['Lon'] for e in self.test_events))
            self.assertLess(cluster['centroid_lon'], max(e['Lon'] for e in self.test_events))

    def test_cluster_hierarchy(self):
        # ทดสอบลำดับชั้นของ clusters
        clusters = build_divisive_tree(self.test_events, max_level=2)
        
        # ตรวจสอบว่า clusters เรียงตาม level
        levels = [c['level'] for c in clusters]
        self.assertEqual(levels, sorted(levels))
        
        # ตรวจสอบว่า parent cluster มี level น้อยกว่า child
        for cluster in clusters:
            if cluster.get('parent_cluster_id'):
                parent = next(c for c in clusters if c['cluster_id'] == cluster['parent_cluster_id'])
                self.assertLess(parent['level'], cluster['level'])

if __name__ == '__main__':
    unittest.main() 