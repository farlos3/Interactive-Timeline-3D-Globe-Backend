from .closest_pair import merge_sort, find_weighted_dist, closest_pair

import numpy as np
from datetime import datetime

def split_cluster_pairs_using_closest_pair(pairs):
    """
    แบ่งกลุ่มโดยใช้ Closest Pair แทนการใช้ radius
    """
    if len(pairs) < 2:
        return pairs, []
    
    # แปลง pairs เป็น list ของ vectors
    vectors = [v for v, _ in pairs]
    
    # ใช้ Closest Pair หาคู่ที่ใกล้ที่สุด
    Px = merge_sort(list(zip(vectors, range(len(vectors)))), 0)
    Py = merge_sort(list(zip(vectors, range(len(vectors)))), 2)
    _, A, B = closest_pair([v for v, _ in Px], [v for v, _ in Py])
    
    # หา index ของ A และ B
    idx_a = -1
    idx_b = -1
    for i, (v, _) in enumerate(pairs):
        if np.array_equal(v, A) and idx_a == -1:
            idx_a = i
        elif np.array_equal(v, B) and idx_b == -1:
            idx_b = i
        if idx_a != -1 and idx_b != -1:
            break
    
    # แบ่งกลุ่มโดยใช้คู่ที่ใกล้ที่สุดเป็น seed
    g1, g2 = [], []
    visited = set()
    
    # เพิ่ม seed เข้าไปในกลุ่ม
    g1.append(pairs[idx_a])
    g2.append(pairs[idx_b])
    visited.update([idx_a, idx_b])
    
    # จัดกลุ่มที่เหลือตามระยะทางที่ใกล้กับ seed
    for i, (v, e) in enumerate(pairs):
        if i in visited:
            continue
            
        # คำนวณระยะทางไปยัง seed ของแต่ละกลุ่ม
        dist_to_g1 = find_weighted_dist(v, pairs[idx_a][0])
        dist_to_g2 = find_weighted_dist(v, pairs[idx_b][0])
        
        # เพิ่มเข้าไปในกลุ่มที่ใกล้กว่า
        if dist_to_g1 < dist_to_g2:
            g1.append(pairs[i])
        else:
            g2.append(pairs[i])
        visited.add(i)
    
    return g1, g2

class ClusterNode:
    def __init__(self, pairs, children=None, label=None):
        self.pairs = pairs  # list of (vec, event)
        self.children = children if children else []
        self.label = label  # ใช้สำหรับระบุชื่อคลัสเตอร์ C1, C2

    @property
    def is_leaf(self):
        return len(self.children) == 0

def divisive_custom_tree(pairs, min_cluster_size, label_counter=None, depth=0):
    """
    Tree ที่แบ่งกลุ่มตามหลักการ Divisive Hierarchical Clustering โดยใช้ Closest Pair
    """
    indent = "  " * depth

    if label_counter is None:
        label_counter = {"count": 1}

    # base case: หยุดแบ่งหากขนาดกลุ่มเล็กเกิน
    if len(pairs) < min_cluster_size:
        label = f"C{label_counter['count']}"
        label_counter['count'] += 1
        print(f"{indent}Leaf: {label}, size = {len(pairs)}")
        return ClusterNode(pairs, label=label)

    # แบ่งกลุ่มด้วย Closest Pair
    g1, g2 = split_cluster_pairs_using_closest_pair(pairs)
    print(f"{indent}Split: size = {len(pairs)}, g1 = {len(g1)}, g2 = {len(g2)}")

    # ถ้าแบ่งไม่สำเร็จ (เช่นก1 หรือ ก2 ว่าง)
    if not g1 or not g2:
        label = f"C{label_counter['count']}"
        label_counter['count'] += 1
        print(f"{indent}Unsplitable: {label}, size = {len(pairs)}")
        return ClusterNode(pairs, label=label)

    # สร้างโหนดลูกแบบ recursive
    child1 = divisive_custom_tree(g1, min_cluster_size, label_counter, depth + 1)
    child2 = divisive_custom_tree(g2, min_cluster_size, label_counter, depth + 1)

    return ClusterNode(pairs, children=[child1, child2])

def get_leaf_clusters(node):
    if node.is_leaf:
        return [node]
    leaves = []
    for child in node.children:
        leaves.extend(get_leaf_clusters(child))
    return leaves

def get_leaf_clusters_with_level(node, level=0):
    """
    Get leaf clusters with their level in the tree
    Returns list of tuples (node, level)
    """
    if node.is_leaf:
        return [(node, level)]
    leaves = []
    for child in node.children:
        leaves.extend(get_leaf_clusters_with_level(child, level + 1))
    return leaves

def leaf_to_cluster_dict(node):
    """
    Convert leaf nodes to a dictionary mapping cluster labels to their data
    """
    result = {}
    for leaf in get_leaf_clusters(node):
        result[leaf.label] = leaf.pairs
    return result

def calculate_centroid(pairs):
    lats = []
    lons = []
    for _, event in pairs:
        # lat/lon robust
        lat = event.get('lat', event.get('Lat'))
        lon = event.get('lon', event.get('Lon'))
        if lat is not None and lon is not None:
            try:
                lats.append(float(lat))
                lons.append(float(lon))
            except Exception:
                continue
        # date robust
        date_val = event.get('date', event.get('Date'))
    centroid_lat = float(np.mean(lats)) if lats else 0.0
    centroid_lon = float(np.mean(lons)) if lons else 0.0
    return centroid_lat, centroid_lon

def calculate_bounding_box(events):
    lats = []
    lons = []
    for event in events:
        lat = event.get('lat', event.get('Lat'))
        lon = event.get('lon', event.get('Lon'))
        try:
            lats.append(float(lat))
            lons.append(float(lon))
        except Exception:
            continue
    if not lats or not lons:
        return ""
    min_lat, max_lat = min(lats), max(lats)
    min_lon, max_lon = min(lons), max(lons)
    # WKT POLYGON (lon lat)
    return f"POLYGON(({min_lon} {min_lat}, {min_lon} {max_lat}, {max_lon} {max_lat}, {max_lon} {min_lat}, {min_lon} {min_lat}))"

def calculate_centroid_days(pairs, min_date):
    days = []
    for _, event in pairs:
        date_val = event.get('date', event.get('Date'))
        if date_val is not None:
            if hasattr(date_val, 'timestamp'):
                dt = date_val
            else:
                try:
                    dt = datetime.strptime(str(date_val), "%Y-%m-%d")
                except Exception:
                    continue
            days.append((dt - min_date).days)
    centroid_days = float(np.mean(days)) if days else 0.0
    return centroid_days

def cluster_tree_to_dict(node, parent_id=None, level=0, cluster_list=None, id_counter=None, min_date=None):
    """
    แปลงต้นไม้คลัสเตอร์เป็น list ของ dict ตาม Model Cluster ของ Go
    เพิ่ม centroid_time_days (จำนวนวันเฉลี่ยจาก min_date) เป็น string
    """
    if cluster_list is None:
        cluster_list = []
    if id_counter is None:
        id_counter = {"id": 1}
    if min_date is None:
        # หา min_date จาก event ใน node.pairs
        min_date_candidates = []
        for _, event in node.pairs:
            date_val = event.get('date', event.get('Date'))
            if date_val is not None:
                if hasattr(date_val, 'timestamp'):
                    min_date_candidates.append(date_val)
                else:
                    try:
                        min_date_candidates.append(datetime.strptime(str(date_val), "%Y-%m-%d"))
                    except Exception:
                        continue
        min_date = min(min_date_candidates) if min_date_candidates else datetime(1970,1,1)
    cluster_id = id_counter["id"]
    id_counter["id"] += 1

    centroid_lat, centroid_lon = calculate_centroid(node.pairs)
    centroid_time_days = str(calculate_centroid_days(node.pairs, min_date))
    group_tag = ""
    bounding_box = calculate_bounding_box([event for _, event in node.pairs])
    event_ids = [event.get('EventID', event.get('event_id')) for _, event in node.pairs]

    cluster_list.append({
        "cluster_id": cluster_id,
        "parent_cluster_id": parent_id,
        "centroid_lat": centroid_lat,
        "centroid_lon": centroid_lon,
        "centroid_time_days": centroid_time_days,
        "level": level,
        "group_tag": group_tag,
        "bounding_box": bounding_box,
        "event_ids": event_ids
    })

    for child in getattr(node, "children", []):
        cluster_tree_to_dict(child, cluster_id, level+1, cluster_list, id_counter, min_date)
    return cluster_list