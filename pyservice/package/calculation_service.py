from calculate import (
    create_feature_vector,
    greedy_closest_pair_kdtree_grouping,
    divisive_custom_tree,
    get_leaf_clusters,
    cluster_tree_to_dict
)
from typing import List, Dict

class CalculationService:
    def __init__(self, k: int = 5, min_cluster_size: int = 10):
        self.k = k
        self.min_cluster_size = min_cluster_size

    def process_events(self, events: List[Dict]) -> Dict:
        # ปรับค่า k ตามจำนวน event (กัน error ถ้า k > จำนวน event)
        actual_k = min(self.k, len(events))
        if actual_k != self.k:
            print(f"Adjusting k from {self.k} to {actual_k} due to limited number of events ({len(events)})")
        
        # ขั้นตอนที่ 1: จัดกลุ่มด้วย greedy KD-Tree
        groups = greedy_closest_pair_kdtree_grouping(events, actual_k)

        all_leaf_nodes = []
        all_cluster_dicts = []
        id_counter = {"id": 1}  # ตัวนับสำหรับ cluster_id

        # ขั้นตอนที่ 2: ทำ hierarchical clustering แยกแต่ละกลุ่ม
        for group_idx, group in enumerate(groups):
            min_date = min(event['Date'] for event in group)
            pairs = [
                (create_feature_vector(event['Lat'], event['Lon'], event['Date'], min_date), event)
                for event in group
            ]
            root = divisive_custom_tree(pairs, self.min_cluster_size)
            leaf_nodes = get_leaf_clusters(root)
            all_leaf_nodes.extend(leaf_nodes)

            # สร้าง dictionary สำหรับ export
            cluster_dicts = cluster_tree_to_dict(root, id_counter=id_counter)
            all_cluster_dicts.extend(cluster_dicts)

        # ตรวจสอบความครบถ้วน
        total = sum(len(leaf.pairs) for leaf in all_leaf_nodes)
        is_complete = total == len(events)

        # สร้าง mapping จาก EventID ไปยัง cluster_id
        event_clusters = {}
        for cluster in all_cluster_dicts:
            cluster_id = cluster['cluster_id']
            for leaf in all_leaf_nodes:
                if hasattr(leaf, 'label') and leaf.label == f"C{cluster_id}":
                    for p in leaf.pairs:
                        event = p[1]
                        event_clusters[event['EventID']] = cluster_id

        result = {
            "total_events": len(events),
            "total_clusters": len(all_cluster_dicts),
            "is_complete": is_complete,
            "missing_events": len(events) - total if not is_complete else 0,
            "clusters": all_cluster_dicts,
            "event_clusters": event_clusters
        }

        return result