from calculate import (
    create_feature_vector,
    greedy_closest_pair_kdtree_grouping,
    divisive_custom_tree,
    get_leaf_clusters
)
from typing import List, Dict

class CalculationService:
    def __init__(self, k: int = 5, min_cluster_size: int = 10):
        self.k = k
        self.min_cluster_size = min_cluster_size

    def process_events(self, events: List[Dict]) -> Dict:
        # ปรับค่า k ตามจำนวน event เด่วมาเอาออก
        actual_k = min(self.k, len(events))
        if actual_k != self.k:
            print(f"Adjusting k from {self.k} to {actual_k} due to limited number of events ({len(events)})")
        
        groups = greedy_closest_pair_kdtree_grouping(events, actual_k)

        all_leaf_nodes = []
        
        for group_idx, group in enumerate(groups):
            min_date = min(event['Date'] for event in group)
            pairs = [
                (create_feature_vector(event['Lat'], event['Lon'], event['Date'], min_date), event)
                for event in group
            ]
            root = divisive_custom_tree(pairs, self.min_cluster_size)
            leaf_nodes = get_leaf_clusters(root)
            all_leaf_nodes.extend(leaf_nodes)

        total = sum(len(leaf.pairs) for leaf in all_leaf_nodes)
        is_complete = total == len(events)

        # สร้าง dictionary เก็บ cluster_id ของแต่ละ event
        event_clusters = {}
        for cluster_idx, leaf in enumerate(all_leaf_nodes):
            for pair in leaf.pairs:
                event = pair[1]
                event_clusters[event['EventID']] = cluster_idx

        # สร้างผลลัพธ์
        result = {
            "total_events": len(events),
            "total_clusters": len(all_leaf_nodes),
            "is_complete": is_complete,
            "missing_events": len(events) - total if not is_complete else 0,
            "clusters": [
                {
                    "cluster_id": idx,
                    "events": [pair[1] for pair in leaf.pairs],
                    "size": len(leaf.pairs)
                }
                for idx, leaf in enumerate(all_leaf_nodes)
            ],
            "event_clusters": event_clusters  # เพิ่มข้อมูล cluster_id ของแต่ละ event
        }

        return result 