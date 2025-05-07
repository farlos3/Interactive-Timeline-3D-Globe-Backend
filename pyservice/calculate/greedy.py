from .feature import create_feature_vector
from scipy.spatial import KDTree
from .closest_pair import merge_sort, find_weighted_dist, closest_pair

import numpy as np
import logging

# ตั้งค่าระดับ logging (DEBUG, INFO, WARNING, ERROR)
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def greedy_closest_pair_kdtree_grouping(events, k):
    # logging.info("Start grouping with k = %d", k)

    # สร้าง mapping ระหว่าง index และ event id
    event_id_map = {i: event['EventID'] for i, event in enumerate(events)}
    
    lat_lon_date = [(event['Lat'], event['Lon'], event['Date']) for event in events]
    lat, lon, date_str = zip(*lat_lon_date)

    min_date = min(date_str)
    max_date = max(date_str)
    # logging.info("Datetime range: %s to %s", min_date, max_date)
    
    points_3d = [create_feature_vector(event['Lat'], event['Lon'], event['Date'], min_date) for event in events]
    
    tree = KDTree(points_3d)

    n = len(events)
    visited = [False] * n
    groups = [[] for _ in range(k)]
    seeds = []

    # --- Step 1: เลือก seed เริ่มต้นโดยใช้ Closest Pair ---
    # logging.info("Selecting initial seed pair using Closest Pair")
    Px = merge_sort(list(zip(points_3d, range(n))), 0)
    Py = merge_sort(list(zip(points_3d, range(n))), 2)
    d, A, B = closest_pair([p for p, _ in Px], [p for p, _ in Py])

    # แก้ไขการค้นหา index ของ A และ B
    idx_a = -1
    idx_b = -1
    for p, idx in Px:
        if np.array_equal(p, A) and idx_a == -1:
            idx_a = idx
        elif np.array_equal(p, B) and idx_b == -1:
            idx_b = idx
        if idx_a != -1 and idx_b != -1:
            break
    
    # ตรวจสอบว่าไม่ใช้ seed ซ้ำ
    if idx_a != -1 and idx_b != -1 and idx_a != idx_b:
        seeds.extend([idx_a, idx_b])
    else:
        # ถ้าหา seed ไม่ได้หรือได้ seed ซ้ำ ให้เลือกจุดแรกและจุดที่ไกลที่สุดจากจุดแรก
        idx_a = 0
        max_dist = -1
        idx_b = 1
        for i in range(1, n):
            dist = find_weighted_dist(points_3d[0], points_3d[i])
            if dist > max_dist:
                max_dist = dist
                idx_b = i
        seeds.extend([idx_a, idx_b])
    
    used_idxs = set(seeds)

    # --- Step 2: เลือก seeds ที่เหลือด้วย Farthest Point Sampling ---
    while len(seeds) < k:
        max_dist = -1
        best_idx = -1
        for idx in range(n):
            if idx in used_idxs:
                continue
            dist = min(find_weighted_dist(points_3d[idx], points_3d[s]) for s in seeds)
            if dist > max_dist:
                max_dist = dist
                best_idx = idx
        seeds.append(best_idx)
        used_idxs.add(best_idx)

    # logging.info("Seeds selected: %s", [event_id_map[idx] for idx in seeds])

    # --- Step 3: จัดกลุ่มตาม seeds ด้วย Greedy KDTree Query ---
    for group_idx, seed_idx in enumerate(seeds):
        groups[group_idx].append(events[seed_idx])
        visited[seed_idx] = True

    remaining_idxs = [i for i in range(n) if not visited[i]]
    target_group_size = n // k

    for group_idx, seed_idx in enumerate(seeds):
        # logging.info("Expanding group %d (seed event id: %d)", group_idx, event_id_map[seed_idx])
        while len(groups[group_idx]) < target_group_size and remaining_idxs:
            # ค้นหา point ที่ใกล้ที่สุดกับ seed ของกลุ่มนี้
            dists, idxs = tree.query(points_3d[seed_idx], k=n)
            for idx in idxs:
                if not visited[idx]:
                    groups[group_idx].append(events[idx])
                    visited[idx] = True
                    remaining_idxs.remove(idx)
                    logging.debug("Added event id %d to group %d", event_id_map[idx], group_idx)
                    break
            else:
                break

    # --- Step 4: กระจาย leftover ไปยังกลุ่มที่เล็กที่สุด ---
    for idx in remaining_idxs:
        smallest_group_idx = min(range(k), key=lambda g: len(groups[g]))
        groups[smallest_group_idx].append(events[idx])
        visited[idx] = True
        # logging.debug("Added leftover event id %d to group %d", event_id_map[idx], smallest_group_idx)

    # logging.info("Grouping completed.")
    return groups