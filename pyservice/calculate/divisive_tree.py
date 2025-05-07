from .closest_pair import merge_sort, find_weighted_dist, closest_pair

import numpy as np

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